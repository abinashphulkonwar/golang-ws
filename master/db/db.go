package db

import (
	"log"

	"github.com/zippoxer/bow"
)

type None struct {
	IP     string
	NAME   string
	STATUS string
}

type Connection struct {
	IP string
}

func OpenDB() {
	db, err := bow.Open("test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Create(&None{IP: "", NAME: "", STATUS: ""})
}
