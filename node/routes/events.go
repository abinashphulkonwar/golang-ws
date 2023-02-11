package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/abinashphulkonwar/ws/db"
)

func WsEvents(w http.ResponseWriter, r *http.Request) {
	mt, err := strconv.Atoi(r.URL.Query().Get("mt"))
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("invalid mt"))

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

	println(string(body))

	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("body not found"))
		return
	}

	connction, isNil := db.Connections[data.SendTo]

	if !isNil {
		w.Write([]byte("connection not found"))
		return
	} else {
		respons, errJson := json.MarshalIndent(&data, "", "\t")
		if errJson != nil {
			w.WriteHeader(500)
			w.Write([]byte("body not found"))
			return
		}
		println(string(respons), mt)
		err := connction.C.WriteMessage(mt, respons)
		if err != nil {
			log.Println("write:", err)
		} else {
			println("event sent")
		}
	}
	w.Write([]byte("event sent"))
}
