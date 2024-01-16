package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Permite todas las conexiones desde cualquier origen
			return true
		},
	}
	clients = make(map[*websocket.Conn]*Client)
	rooms   = make(map[string]*Room)
	mu      sync.Mutex
)

// Client representa a un cliente conectado al servidor WebSocket
type Client struct {
	ID         string
	Connection *websocket.Conn
	Room       *Room // Agrega el campo Room
}

// Room representa una sala en el servidor WebSocket
type Room struct {
	ID      string
	Clients map[*Client]bool
	mu      sync.Mutex
}

func handleWebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// ID del cliente
	clientID := c.QueryParam("id")
	if clientID == "" {
		clientID = generateUniqueID()
	}
	// Sala de chat
	room := c.QueryParam("room")
	if room == "" {
		room = "public"
	}

	client := &Client{
		ID:         clientID,
		Connection: conn,
	}

	// Asocia el cliente con la sala
	joinRoom(client, room)

	// Imprime un mensaje cuando un cliente se conecta
	fmt.Printf("Cliente %s conectado al WebSocket\n", client.ID)

	for {

		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error al leer el mensaje:", err)
			break
		}

		fmt.Printf("Mensaje recibido de %s en la sala %s: %s\n", client.ID, client.Room.ID, msg)

		broadcast(client.Room, msg)
	}

	// Elimina al cliente de la sala y del mapa cuando se desconecta
	leaveRoom(client)
	fmt.Printf("Cliente %s desconectado del WebSocket\n", client.ID)

	return nil
}

func joinRoom(client *Client, roomID string) {
	mu.Lock()
	defer mu.Unlock()

	// Si la sala no existe se crea
	if _, ok := rooms[roomID]; !ok {
		rooms[roomID] = &Room{
			ID:      roomID,
			Clients: make(map[*Client]bool),
		}
	}

	// Asocia al cliente con la sala
	client.Room = rooms[roomID]
	client.Room.Clients[client] = true
}

/** Elimina al cliente de la sala y del mapa cuando se desconecta*/
func leaveRoom(client *Client) {
	mu.Lock()
	defer mu.Unlock()

	if client.Room != nil {
		// Elimina al cliente de la sala
		delete(client.Room.Clients, client)

		// Si la sala esta vacía se elimina
		if len(client.Room.Clients) == 0 {
			delete(rooms, client.Room.ID)
		}
	}
}

/** Envía el mensaje a todos los clientes en la misma sala */
func broadcast(room *Room, message []byte) {
	room.mu.Lock()
	defer room.mu.Unlock()

	for client := range room.Clients {
		err := client.Connection.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Printf("Error al enviar mensaje a %s: %s\n", client.ID, err)
		}
	}
}

/** Genera un ID único para cada cliente conectado al servidor WebSocket  */
func generateUniqueID() string {
	return fmt.Sprintf("Client%d", len(clients)+1)
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/ws", handleWebSocket)
	e.Logger.Fatal(e.Start(":8080"))
}
