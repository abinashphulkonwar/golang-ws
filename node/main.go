package main

import (
	"encoding/json"
	"flag"
	"io"

	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/abinashphulkonwar/ws/db"
	"github.com/abinashphulkonwar/ws/routes"
	"github.com/abinashphulkonwar/ws/service"
	"github.com/abinashphulkonwar/ws/src"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request, ip string) {

	println("echo", r.Host, r.URL, r.URL.Path)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()

	id := src.Uuid()
	connection := db.Connection{Id: id, C: c, Node: ip, Ttl: time.Now()}
	db.Connections[id] = &connection
	println(r.Cookies())
	println(id)
	service.SetConnections(&service.Connection{Id: id, Node: ip})

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			err := db.ErrorRes{Message: "error while parsing message: " + err.Error(), Type: db.ErrorType}
			println(err.Message)
			errMessage, _ := json.MarshalIndent(&err, "", "\t")
			c.WriteMessage(mt, errMessage)
			return
		}
		chat := db.Chat{}

		errJson := json.Unmarshal(message, &chat)
		if errJson != nil {
			err := db.ErrorRes{Message: "error while parsing message: " + errJson.Error(), Type: db.ErrorType}
			println(err.Message)
			errMessage, _ := json.MarshalIndent(&err, "", "\t")
			c.WriteMessage(mt, errMessage)
			return
		}

		if chat.SendTo == "" {
			return
		}
		chat.SendFrom = id
		chat.Type = db.MessageType
		respons, err := json.MarshalIndent(&chat, "", "\t")
		if err != nil {
			err := db.ErrorRes{Message: "error while parsing message: " + errJson.Error(), Type: db.ErrorType}
			errMessage, _ := json.MarshalIndent(&err, "", "\t")
			c.WriteMessage(mt, errMessage)
			return
		}

		log.Printf("recv: %s", message)
		// err = c.WriteMessage(mt, respons)
		// if err != nil {
		// 	err := db.ErrorRes{Message: "error : " + err.Error(), Type: db.ErrorType}
		// 	println(err.Message)
		// 	errMessage, _ := json.MarshalIndent(&err, "", "\t")
		// 	c.WriteMessage(mt, errMessage)
		// 	return
		// }
		println(chat.SendTo, mt)
		ws, isNil := db.Connections[chat.SendTo]

		if !isNil {

			node, isNil := db.Connectionothers[chat.SendTo]
			exp := node.Ttl.Add(25 * time.Second)
			if isNil && exp.Before(time.Now()) {
				service.PostEvent(&db.ConnectionRes{
					Message: "node not found",
					Node:    node.Node,
				}, mt, respons)

			} else {
				res, err := service.GetNode(&chat)
				if err != nil {
					println("error while fetching")
				} else {
					data, _ := io.ReadAll(res)
					println(string(data))
					body := db.ConnectionRes{}
					err = json.Unmarshal(data, &body)
					if err != nil {
						println("error while parsing")

					} else {

						db.Connectionothers[chat.SendTo] = &db.OthersConnection{
							Id:   chat.SendTo,
							Node: body.Node,
							Ttl:  time.Now(),
						}

						service.PostEvent(&body, mt, respons)
					}

				}
			}

		} else {
			println(chat.Message)
			if ws.Id != id {
				err := ws.C.WriteMessage(mt, respons)
				if err != nil {
					log.Println("write:", err)
				}
			}
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	println("home", r.Host, r.URL, r.URL.Path)
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {

	port := "3001"

	db.UpdateTtl()

	for i := 1; i < len(os.Args); i++ {
		println(os.Args[i])
		if os.Args[i] == "p" || os.Args[i] == "port" {
			if os.Args[i+1] != "" {
				port = os.Args[i+1]
				break
			} else {
				println("port not provided")

			}

		}
	}
	ip := "localhost:" + port
	var addr = flag.String("addr", ip, "http service address")

	flag.Parse()
	log.SetFlags(0)
	go http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		echo(w, r, ip)
	})
	http.HandleFunc("/", home)
	http.HandleFunc("/events", routes.WsEvents)

	service.SetNode(&service.Node{IP: ip,
		NAME:   port,
		STATUS: "active"})

	println("listening on port", port)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
			const data = JSON.parse(evt.data)
			
			if (data.Type === "message") {
          print("RESPONSE: " + data.message + " from " + data.from);
			}

			if (data.Type === "error") {
              print("Error: " + data.message);
			}

        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
	
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
		const userId = document.getElementById("user").value
        print("SEND: " + input.value);
       ws.send(JSON.stringify({id:userId , message: input.value}));
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="user" type="text" value="user">
<p><input id="input" type="text" value="">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
