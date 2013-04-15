package dispatch

import (
	"log"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/nickpresta/gonotify/mailbox"
)

// A connection to a client, a channel for sending a message, and a Mailbox
type Connection struct {
	WebSocket   *websocket.Conn
	Send chan string
	Mailbox *mailbox.Mailbox
}

type jsonMessage struct {
	Date    string `json:"date"`
	Message string `json:"message"`
}

const (
	pingPeriod = 30 * time.Second // How often to check for dead connections
)

// Checks if the connection is still alive
// Sends a PING frame over the websocket
func (c *Connection) isAlive() bool {
	// We need to change the payload type to a PING frame
	// Save the old payload type, change it to PING, then restore
	oldPayloadType := c.WebSocket.PayloadType
	c.WebSocket.PayloadType = websocket.PingFrame
	n, err := c.WebSocket.Write([]byte{})
	c.WebSocket.PayloadType = oldPayloadType
	return n == 0 && err == nil
}

// Writes messages to the connection, checks for dead connections.
func (c *Connection) Writer(dispatcher Dispatch) {
	ping := time.NewTicker(pingPeriod)
WriteLoop:
	for {
		select {
		case message := <-c.Send:
			now := time.Now().Format(time.UnixDate)
			jsonMessage := jsonMessage{
				Date:    now,
				Message: message,
			}
			if err := websocket.JSON.Send(c.WebSocket, jsonMessage); err != nil {
				break
			}
		case <-ping.C:
			log.Printf("Pinging: %v (Connection: %p)\n", c.Mailbox.Receiver, c)
			if !c.isAlive() {
				break WriteLoop
			}
		}
	}
	ping.Stop()
	dispatcher.Unregister <- c
}
