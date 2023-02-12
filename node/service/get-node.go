package service

import (
	"errors"
	"io"
	"strconv"

	"github.com/abinashphulkonwar/ws/db"
)

func GetNode(c *db.Chat) (io.Reader, error) {

	res, err := Fetch(&Request{
		Method: "GET",
		Url:    "http://localhost:3000/connections/get/?id=" + c.SendTo,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	if err != nil && res.StatusCode != 200 {
		println("error while fetching")
		return nil, errors.New("error while fetching")
	} else {

		return res.Body, nil

	}

}

func PostEvent(body *db.ConnectionRes, mt int, respons []byte) {
	resMessage, _ := Fetch(&Request{
		Method: "POST",
		Url:    "http://" + body.Node + "/events?" + "mt=" + strconv.Itoa(mt),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: respons,
	})
	println(resMessage.StatusCode)
}
