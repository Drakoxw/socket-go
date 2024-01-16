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
	e.GET("/ws", chat.HandleWebSocket)

	port := utils.GetPort()
	e.Logger.Fatal(e.Start(port))
}
