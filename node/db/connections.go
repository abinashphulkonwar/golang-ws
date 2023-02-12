package db

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/jasonlvhit/gocron"
)

type Connection struct {
	Id   string
	C    *websocket.Conn
	Node string
	Ttl  time.Time
}
type OthersConnection struct {
	Id   string
	Node string
	Ttl  time.Time
}

type ConnectionRes struct {
	Message string `json:"message"`
	Node    string `json:"node"`
}

var Connections = make(map[string]*Connection)
var Connectionothers = make(map[string]*OthersConnection)

func UpdateTtl() {
	s := gocron.NewScheduler()

	s.Every(5).Minutes().Do(func() {
		currentTime := time.Now()

		for key, value := range Connectionothers {

			exp := value.Ttl.Add(25 * time.Second)

			if currentTime.After(exp) {
				delete(Connectionothers, key)
			}
		}

	})
}
