package dispatch

import (
	"log"

	"github.com/nickpresta/gonotify/mailbox"
)

type connectionArray []*Connection

type Dispatch struct {
	connections map[*Connection]bool
	nameConnectionsLookup map[string]connectionArray
	Register    chan *Connection
	Unregister  chan *Connection
	Broadcast   chan *mailbox.Mailbox
}

// Dispatcher that provides register, unregister, and broadcast channels
func NewDispatcher() Dispatch {
	return Dispatch{
		connections: make(map[*Connection]bool),
		nameConnectionsLookup:  make(map[string]connectionArray),
		Register:    make(chan *Connection),
		Unregister:  make(chan *Connection),
		Broadcast:   make(chan *mailbox.Mailbox),
	}
}

func (c connectionArray) find(value *Connection) int {
	for i, conn := range(c) {
		if conn == value {
			return i
		}
	}
	return -1
}

func (d *Dispatch) Run() {
	for {
		select {
		case conn := <-d.Register:
			log.Printf("Registering: %v (Connection: %p)\n", conn.Mailbox.Receiver, conn)
			d.connections[conn] = true
			d.nameConnectionsLookup[conn.Mailbox.Receiver] = append(d.nameConnectionsLookup[conn.Mailbox.Receiver], conn)
			log.Printf("Number of connections: %d\n", len(d.connections))
		case conn := <-d.Unregister:
			log.Printf("Unregistering: %v (Connection: %p)\n", conn.Mailbox.Receiver, conn)
			delete(d.connections, conn)

			// Delete the connection from the slice for that receiver
			// We have to delete the data this way, otherwise we risk a memory leak (as it won't be GC'd)
			connSlice := d.nameConnectionsLookup[conn.Mailbox.Receiver]
			removalIndex := connSlice.find(conn)
			connSliceLen := len(connSlice)
			copy(connSlice[removalIndex:], connSlice[removalIndex+1:])
			connSlice[connSliceLen-1] = nil
			d.nameConnectionsLookup[conn.Mailbox.Receiver] = connSlice[:connSliceLen-1]
			if len(d.nameConnectionsLookup[conn.Mailbox.Receiver]) == 0 {
				delete(d.nameConnectionsLookup, conn.Mailbox.Receiver)
			}

			log.Printf("Number of connections: %d\n", len(d.connections))
			close(conn.Send)
			conn.WebSocket.Close()
		case mail := <-d.Broadcast:
			log.Printf("Got broadcast %+v\n", mail)
			if connections, ok := d.nameConnectionsLookup[mail.Receiver]; ok {
				for _, conn := range(connections) {
					log.Printf("Sending message to: %v (Connection: %p)\n", mail.Receiver, conn)
					conn.Send <- mail.Message
				}
			}
		}
	}
}
