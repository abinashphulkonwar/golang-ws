package service

import (
	"bytes"
	"errors"
	"net/http"
)

type Request struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    string
}

var client = &http.Client{}

func req(Request *Request) (*http.Request, error) {

	if Request.Method == "GET" {
		req, err := http.NewRequest(Request.Method, Request.Url, nil)
		return req, err
	}
	if Request.Body != "" {

		data := bytes.NewBuffer([]byte(Request.Body))

		return http.NewRequest(Request.Method, Request.Url, data)
	}
	req, err := http.NewRequest(Request.Method, Request.Url, nil)
	return req, err

}

func Fetch(Request *Request) (*http.Response, error) {

	if Request.Method == "GET" {
		if Request.Body != "" {
			return nil, errors.New("body is not allowed for get request")
		}
	}
	req, err := req(Request)

	if err != nil {
		return nil, err
	}
	for key, value := range Request.Headers {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
