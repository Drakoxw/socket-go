// rest_proxy.go

package chat

import (
	"fmt"
	"net/http"
	"time"

	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// RestProxy es una estructura que maneja las peticiones REST y las envía al socket
type RestProxy struct {
	Upgrader websocket.Upgrader
	Clients  map[*websocket.Conn]*Client
	Rooms    map[string]*Room
	Mu       sync.Mutex
}


func (rp *RestProxy) HandleRestRequest(c echo.Context) error {

	var mensaje Message
	if err := c.Bind(&mensaje); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Formato de mensaje JSON inválido"})
	}
	mensaje.DateTime = time.Now()

	room := c.QueryParam("room")
	if room == "" {
		room = "public"
	}

	// Envia el mensaje al socket
	rp.Mu.Lock()
	defer rp.Mu.Unlock()

	// Si la sala no existe, crea una nueva
	if _, ok := rp.Rooms[room]; !ok {
		rp.Rooms[room] = &Room{
			ID:      room,
			Clients: make(map[*Client]bool),
		}
	}

	fmt.Printf("Rooms: %+v\n", rp.Rooms)

	// Envía el mensaje a todos los clientes en la sala
	for client := range rp.Rooms[room].Clients {
		err := client.Connection.WriteJSON(mensaje)
		if err != nil {
			fmt.Printf("Error al enviar mensaje a %s: %s\n", client.ID, err)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Mensaje enviado al socket con éxito"})
}

var Proxy = &RestProxy{
	Upgrader: upgrader,
	Clients:  clients,
	Rooms:    rooms,
	Mu:       mu,
}
