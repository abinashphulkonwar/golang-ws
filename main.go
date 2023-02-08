package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/abinashphulkonwar/ws/routes"
	"github.com/abinashphulkonwar/ws/src"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Connection struct {
	id string
	c  *websocket.Conn
}

var connections = make(map[string]*Connection)

var ids = []string{}

func echo(w http.ResponseWriter, r *http.Request) {
	println("echo", r.Host, r.URL, r.URL.Path)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()

	id := src.Uuid()
	ids = append(ids, id)
	connection := Connection{id: id, c: c}
	connections[id] = &connection
	println(r.Cookies())
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		res := append([]byte("server  "), message...)
		err = c.WriteMessage(mt, res)
		for i := 0; i < len(ids); i++ {
			ws := *connections[ids[i]]
			if ws.id != id {
				err := ws.c.WriteMessage(mt, append([]byte(ws.id+" "), message...))
				if err != nil {
					log.Println("write:", err)
				}
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

	port := "3000"

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
	var addr = flag.String("addr", "localhost:"+port, "http service address")

	flag.Parse()
	log.SetFlags(0)
	go http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	http.HandleFunc("/events", routes.WsEvents)
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
