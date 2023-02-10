package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/abinashphulkonwar/ws/db"
)

func WsEvents(w http.ResponseWriter, r *http.Request) {

	mt, err := strconv.Atoi(r.URL.Query().Get("mt"))
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("method not allowed"))

		return
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Write([]byte("method not allowed"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("body not found"))
		return
	}

	data := db.Chat{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("body not found"))
		return
	}

	connction, isNil := db.Connections[data.SendTo]

	if isNil {
		w.WriteHeader(405)
		w.Write([]byte("connection not found"))
		return
	}

	connction.C.WriteMessage(mt, body)

	w.Write([]byte("event sent"))
}
