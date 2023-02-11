package db

const MessageType = "message"
const ErrorType = "error"

type Chat struct {
	SendTo   string `json:"id"`
	SendFrom string `json:"from"`
	Message  string `json:"message"`
	Type     string `json:"Type"`
}

type ErrorRes struct {
	Message  string `json:"message"`
	SendFrom string `json:"from"`
	Type     string `json:"Type"`
}
