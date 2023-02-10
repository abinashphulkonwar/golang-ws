package main

import (
	"encoding/json"
	"flag"
	"strconv"

	"html/template"
	"io/ioutil"
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

type ConnectionRes struct {
	Message string `json:"message"`
	Node    string `json:"node"`
}

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
			log.Println("read:", err)
			break
		}
		chat := db.Chat{}
		errJson := json.Unmarshal(message, &chat)
		if errJson != nil {
			c.WriteMessage(mt, []byte("error while parsing message"))
			return
		}
		log.Printf("recv: %s", message)
		res := append([]byte(id+"   "), message...)
		err = c.WriteMessage(mt, res)

		ws, isNil := db.Connections[chat.SendTo]
		if isNil {
			res, err := service.Fetch(&service.Request{
				Method: "GET",
				Url:    "http://localhost:3000/connections/get/?id=" + chat.SendTo,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			})

			if err != nil {
				println("error while fetching")
				return
			}
			if res.StatusCode != 200 {
				println("error while fetching")
				return
			}
			data, _ := ioutil.ReadAll(res.Body)

			println(string(data))
			body := ConnectionRes{}
			err = json.Unmarshal(data, &body)
			if err != nil {
				println("error while parsing")
				return
			}

			resMessage, error := service.Fetch(&service.Request{
				Method: "GET",
				Url:    body.Node + "events?" + "mt=" + strconv.Itoa(mt),
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: message,
			})

			if error != nil {
				println("error while fetching")
				return
			}
			if resMessage.StatusCode != 200 {
				println("error while fetching")
				return
			}

			return

		}

		if ws.Id != id {
			err := ws.C.WriteMessage(mt, append([]byte(id+" "), message...))
			if err != nil {
				log.Println("write:", err)
			}
		}

		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	println("home", r.Host, r.URL, r.URL.Path)
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {

	port := "3001"

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
		NAME:   "node1",
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
          print("RESPONSE: " + evt.data);
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
        print("SEND: " + input.value);
        ws.send(input.value);
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
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
