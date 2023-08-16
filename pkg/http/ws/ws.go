package ws

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Use One global Hub
var hub *Hub

// serveWs handles websocket requests from the peer.
func ServeWs(ginContext *gin.Context) {
	conn, err := upgrader.Upgrade(ginContext.Writer, ginContext.Request, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// Broadcast data to the global hub
func Broadcast(data []byte) {
	hub.broadcast <- data
}

func init() {
	// Init our global hub
	hub = newHub()

	go hub.run()
}
