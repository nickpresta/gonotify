package dispatch

import (
	"log"

	"github.com/nickpresta/gonotify/mailbox"
)

type Dispatch struct {
	connections map[string]*Connection
	Register    chan *Connection
	Unregister  chan *Connection
	Broadcast   chan *mailbox.Mailbox
}

// Dispatcher that provides register, unregister, and broadcast channels
func NewDispatcher() Dispatch {
	return Dispatch{
		connections: make(map[string]*Connection),
		Register:    make(chan *Connection),
		Unregister:  make(chan *Connection),
		Broadcast:   make(chan *mailbox.Mailbox),
	}
}

func (d *Dispatch) Run() {
	for {
		select {
		case conn := <-d.Register:
			log.Printf("Registering: %v\n", conn.Mailbox.Receiver)
			d.connections[conn.Mailbox.Receiver] = conn
		case conn := <-d.Unregister:
			log.Printf("Unregistering: %v\n", conn.Mailbox.Receiver)
			delete(d.connections, conn.Mailbox.Receiver)
			log.Printf("Connections: %v\n", d.connections)
			close(conn.Send)
		case mail := <-d.Broadcast:
			log.Printf("Got broadcast %+v\n", mail)
			if conn, ok := d.connections[mail.Receiver]; ok {
				log.Printf("Sending message to: %v\n", mail.Receiver)
				conn.Send <- mail.Message
			}
		}
	}
}
