package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Node struct {
	IP     string `bow:"key"` // primary key
	NAME   string
	STATUS string
}

func SetNode(node *Node) error {
	payload, err := json.Marshal(node)
	if err != nil {
		return errors.New("error while marshaling node")
	}
	res, err := Fetch(&Request{
		Method: "POST",
		Url:    "http://localhost:3000/nodes/post",
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
