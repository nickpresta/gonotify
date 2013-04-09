package main

import (
	"fmt"
)

type dispatcher struct {
	connections map[string]*connection
	register    chan *connection
	unregister  chan *connection
	broadcast   chan *mail
}

var dispatch = dispatcher{
	connections: make(map[string]*connection),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	broadcast:   make(chan *mail),
}

func (d *dispatcher) run() {
	for {
		select {
		case conn := <-d.register:
			d.connections[conn.mail.receiver] = conn
		case conn := <-d.unregister:
			fmt.Println("unregistering")
			delete(d.connections, conn.mail.receiver)
			close(conn.send)
		case mail := <-d.broadcast:
			if conn, ok := d.connections[mail.receiver]; ok {
				conn.send <- mail.message
			}
		}
	}
}
