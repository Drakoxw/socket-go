package main

import (
	"net/http"
	"socket/internal/app/chat"
	"socket/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func HandleHello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hola mundo!"})
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	// Endpoint para manejar las peticiones REST a Socket
	e.POST("/rest-sw", chat.Proxy.HandleRestRequest)

	// WebSocket
	e.GET("/ws", chat.HandleWebSocket)

	e.GET("/", HandleHello)

	port := utils.GetPort()
	e.Logger.Fatal(e.Start(port))
}
