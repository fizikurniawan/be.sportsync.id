package ws

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

var Clients = make(map[*websocket.Conn]string)
var ClientsMu sync.Mutex

func DisconnectClient(conn *websocket.Conn) {
	// Lock the mutex to ensure safe access to the Clients map
	ClientsMu.Lock()
	defer ClientsMu.Unlock()

	// Check if the connection exists in the Clients map
	if _, exists := Clients[conn]; exists {
		// Close the WebSocket connection
		conn.Close()

		// Remove the connection from the Clients map
		delete(Clients, conn)
	}
}
