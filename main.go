package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

type mail struct {
	receiver string
	message  string
}

type JSONPostData struct {
	Mailbox string
	Message string
}

var (
	port          = flag.Int("port", 8080, "HTTP listen port")
	indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
)

func mainRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mailbox := vars["mailbox"]
	context := struct {
		Port    int
		Mailbox string
	}{
		*port,
		mailbox,
	}
	indexTemplate.Execute(w, context)
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

	mail := &mail{receiver: data.Mailbox, message: data.Message}
	dispatch.broadcast <- mail
}

func websocketRequestHandler(ws *websocket.Conn) {
	vars := mux.Vars(ws.Request())
	receiver := vars["mailbox"]

	mail := &mail{receiver: receiver}
	c := &connection{send: make(chan string, 1024), ws: ws, mail: mail}
	dispatch.register <- c
	defer func() { dispatch.unregister <- c }()
	c.writer()
}

func main() {
	flag.Parse()

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	go dispatch.run()

	router := mux.NewRouter()
	router.HandleFunc("/mailbox/{mailbox}", mainRequestHandler).Methods("GET")
	router.HandleFunc("/send", sendRequestHandler).Methods("POST")
	router.Handle("/websocket/mailbox/{mailbox}", websocket.Handler(websocketRequestHandler)).Methods("GET")

	http.Handle("/", router)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
