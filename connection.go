package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"time"
)

// A websocket connection and the output message to send
type connection struct {
	ws   *websocket.Conn
	send chan string
	mail *mail
}

type JSONMessage struct {
	Date    string `json:"date"`
	Message string `json:"message"`
}

func createJSONPayload(message string) string {
	now := time.Now().Format(time.UnixDate)
	jsonMessage := JSONMessage{
		Date:    now,
		Message: message,
	}
	payload, err := json.Marshal(jsonMessage)
	if err != nil {
		return "{}"
	}
	return string(payload)
}

func (c *connection) writer() {
	for {
		select {
		case message := <-c.send:
			payload := createJSONPayload(message)
			if err := websocket.Message.Send(c.ws, payload); err != nil {
				break
			}
		}
	}
	c.ws.Close()
}
