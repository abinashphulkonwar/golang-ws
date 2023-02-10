package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Connection struct {
	Id   string `bow:"key"`
	Node string
}

func SetConnections(Connection *Connection) error {
	payload, err := json.Marshal(Connection)
	if err != nil {
		return errors.New("error while marshaling node")
	}
	res, err := Fetch(&Request{
		Method: "POST",
		Url:    "http://localhost:3000/connections/post",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: payload,
	})

	if err != nil {
		return errors.New("error while fetching")
	}
	body := res.Body

	data, _ := ioutil.ReadAll(body)

	println(string(data))
	if res.StatusCode != 200 {
		return errors.New("error")
	}

	return nil
}
