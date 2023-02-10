package db

import (
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Id   string
	C    *websocket.Conn
	Node string
	Ttl  time.Time
}

var Connections = make(map[string]*Connection)
