package main

import (
	"socket/internal/app/chat"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/ws", chat.HandleWebSocket)
	e.Logger.Fatal(e.Start(":8080"))
}
