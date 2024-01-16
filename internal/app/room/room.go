package room

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Room struct {
	ID      string
	Clients map[*Client]bool
	mu      sync.Mutex
}

type Client struct {
	ID         string
	Connection *websocket.Conn
	Room       *Room
}
