package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/gorilla/mux"
	"github.com/nickpresta/gonotify/dispatch"
)

var (
	port          = flag.Int("port", 8080, "HTTP listen port")
	debugOn       = flag.Bool("debug", false, "Turn on debug logging")
	indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
	dispatcher    = dispatch.NewDispatcher()
)

const (
	statsInterval = 10 * time.Second
)

func mainRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mailbox := vars["mailbox"]
	context := map[string]interface{}{"Port": *port, "Mailbox": mailbox}
	indexTemplate.Execute(w, context)
}

func main() {
	flag.Parse()

	if !*debugOn {
		log.SetOutput(ioutil.Discard)
	}

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	go func() {
		var stats runtime.MemStats
		delay := time.NewTicker(statsInterval)
		defer delay.Stop()
		for {
			select {
			case <-delay.C:
				runtime.ReadMemStats(&stats)
				log.Printf("Alloc: %f", float64(stats.Alloc))
				log.Printf("Sys: %f", float64(stats.Sys))
				log.Printf("Goroutines: %f", float64(runtime.NumGoroutine()))
			}
		}
	}()

	go dispatcher.Run()

	router := mux.NewRouter()
	router.HandleFunc("/mailbox/{mailbox}", mainRequestHandler).Methods("GET")
	router.HandleFunc("/send", sendRequestHandler).Methods("POST")
	router.Handle("/websocket/mailbox/{mailbox}", websocket.Handler(websocketRequestHandler))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/", router)

	log.Printf("Serving on 0.0.0.0:%d...\n", *port)
	if err := http.ListenAndServeTLS(fmt.Sprintf("0.0.0.0:%d", *port), "cert.pem", "key.pem", nil); err != nil {
		log.Fatal(err)
	}
}
