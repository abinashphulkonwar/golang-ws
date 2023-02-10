package service

import (
	"io/ioutil"
	"testing"
)

func TestFetch(t *testing.T) {
	res, err := Fetch(&Request{
		Method: "GET",
		Url:    "http://localhost:3000/",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	println(res.StatusCode, res.Status)

	body := res.Body

	data, _ := ioutil.ReadAll(body)

	println(string(data))
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Error("Status code is not 200")
	}

}
