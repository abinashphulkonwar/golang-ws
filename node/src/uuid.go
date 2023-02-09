package src

import (
	"log"

	"github.com/teris-io/shortid"
)

func Uuid() string {
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		log.Print("upgrade:", err)
		panic("uuid generation")
	}
	id, err := sid.Generate()
	if err != nil {
		panic("uuid error")
	}
	return id
}
