package ws

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

var sessionGroupMap = make(map[string]socketio.Conn)

func CreateWebsocketServer() *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		sessionGroupMap[s.ID()] = s
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		delete(sessionGroupMap, s.ID())
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		delete(sessionGroupMap, s.ID())
		log.Println("closed", reason)
	})

	return server
}

func Broadcast(threadName string, messageContent []byte) {
	for _, wsSession := range sessionGroupMap {
		wsSession.Emit(threadName, messageContent)
	}
}
