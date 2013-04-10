package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"

	"code.google.com/p/go.net/websocket"
	"github.com/nickpresta/gonotify/dispatch"
	"github.com/nickpresta/gonotify/mailbox"
)

type JSONPostData struct {
	Mailbox string
	Message string
}

func sendRequestHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse POST body: %+v", err), http.StatusInternalServerError)
	}
	var data JSONPostData
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse POST body: %+v", err), http.StatusInternalServerError)
	}

	mail := &mailbox.Mailbox{Receiver: data.Mailbox, Message: data.Message}
	dispatcher.Broadcast <- mail
}

func websocketRequestHandler(ws *websocket.Conn) {
	vars := mux.Vars(ws.Request())
	receiver := vars["mailbox"]

	mail := &mailbox.Mailbox{Receiver: receiver}
	c := &dispatch.Connection{Send: make(chan string, 1024), WebSocket: ws, Mailbox: mail}
	dispatcher.Register <- c
	c.Writer(dispatcher)
}
