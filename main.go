package main

import (
	"socket/internal/app/chat"
	"socket/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	// Endpoint para manejar las peticiones REST a Socket
	e.POST("/rest-sw", chat.Proxy.HandleRestRequest)

	// WebSocket
	e.GET("/ws", chat.HandleWebSocket)

	port := utils.GetPort()
	e.Logger.Fatal(e.Start(port))
}
