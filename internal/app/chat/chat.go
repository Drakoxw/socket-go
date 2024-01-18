package chat

import (
	"fmt"
	"socket/pkg/socket_cors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = socket_cors.Upgrader
	clients  = make(map[*websocket.Conn]*Client)
	rooms    = make(map[string]*Room)
	mu       sync.Mutex
)

type Client struct {
	ID         string
	Connection *websocket.Conn
	Room       *Room
}

type Room struct {
	ID      string
	Clients map[*Client]bool
	mu      sync.Mutex
}

type Message struct {
	Sender   int64     `json:"sender"`
	Content  string    `json:"content"`
	DateTime time.Time `json:"dateTime"`
}

func HandleWebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	clientID := c.QueryParam("id")
	if clientID == "" {
		clientID = generateUniqueID()
	}
	room := c.QueryParam("room")
	if room == "" {
		room = "public"
	}

	client := &Client{
		ID:         clientID,
		Connection: conn,
	}
	joinRoom(client, room)
	fmt.Printf("Cliente %s conectado al WebSocket\n", client.ID)

	for {
		var message Message
		err := conn.ReadJSON(&message)
		if err != nil {
			fmt.Println("Error al leer el mensaje:", err)
			break
		}
		fmt.Printf("Mensaje recibido de %s en la sala '%s': %+v\n", client.ID, client.Room.ID, message)
		broadcast(client.Room, message)
	}

	leaveRoom(client)
	fmt.Printf("Cliente %s desconectado del WebSocket\n", client.ID)

	return nil
}

func joinRoom(client *Client, roomID string) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := rooms[roomID]; !ok {
		rooms[roomID] = &Room{
			ID:      roomID,
			Clients: make(map[*Client]bool),
		}
	}

	client.Room = rooms[roomID]
	client.Room.Clients[client] = true
}

func leaveRoom(client *Client) {
	mu.Lock()
	defer mu.Unlock()

	if client.Room != nil {
		delete(client.Room.Clients, client)

		if len(client.Room.Clients) == 0 {
			delete(rooms, client.Room.ID)
		}
	}
}

func broadcast(room *Room, message Message) {
	message.DateTime = time.Now()
	room.mu.Lock()
	defer room.mu.Unlock()

	for client := range room.Clients {
		err := client.Connection.WriteJSON(message)
		if err != nil {
			fmt.Printf("Error al enviar mensaje a %s: %s\n", client.ID, err)
		}
	}
}

func generateUniqueID() string {
	return fmt.Sprintf("Client%d", len(clients)+1)
}
