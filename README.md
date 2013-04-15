# Gonotify
Gonotify is a proof of concept notification server.

![Powered by Gophers](http://i.imgur.com/SwkPj.png "Powered by Gophers")

## How to run
1. Get gonotify: `go get github.com/nickpresta/gonotify`
2. Build gonotify: `go build github.com/nickpresta/gonotify`

You should now have a binary, `gonotify`. See `gonotify --help` for available options.

## How to interact
Gonotify is an HTTP and WebSocket server.
It accepts messages via HTTP and relays them to clients connected via WebSocket.

### Clients
* Browse to `http://localhost:8080/mailbox/your-mailbox-name` - e.g. `http://localhost:8080/mailbox/nickp`
* Wait

### To send data
* Send a `POST` request to `http://localhost:8080/send` with a JSON payload looking like:

        {
            "mailbox": "nickp",
            "message": "hello!"
        }

* The mailbox should receive the message and show up on the client's screen.
