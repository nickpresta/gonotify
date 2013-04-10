package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nickpresta/gonotify/dispatch"
	"html/template"
	"log"
	"net/http"
	"runtime"
)

var (
	port          = flag.Int("port", 8080, "HTTP listen port")
	indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
	dispatcher    = dispatch.NewDispatcher()
)

func mainRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mailbox := vars["mailbox"]
	context := map[string]interface{}{"Port": *port, "Mailbox": mailbox}
	indexTemplate.Execute(w, context)
}

func main() {
	flag.Parse()

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	go dispatcher.Run()

	router := mux.NewRouter()
	router.HandleFunc("/mailbox/{mailbox}", mainRequestHandler).Methods("GET")
	router.HandleFunc("/send", sendRequestHandler).Methods("POST")
	router.Handle("/websocket/mailbox/{mailbox}", websocket.Handler(websocketRequestHandler))

	http.Handle("/", router)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
