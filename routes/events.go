package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type bodyType struct {
	Value string
	Valid bool
}

func WsEvents(w http.ResponseWriter, r *http.Request) {
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

	data := bodyType{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("body not found"))
		return
	}

	println("event", data.Valid, data.Valid)
	w.Write([]byte("event sent"))
}
