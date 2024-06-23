package controllers

import (
	"context"
	"rapidEx/internal/services/auth"
	"rapidEx/internal/services/order"
	"rapidEx/internal/services/stock"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Controllers struct {
	ctx          context.Context
	authService  auth.AuthService
	orderService order.OrderService
	stockService stock.StockService
}

// TODO: normal error handling
func (c *Controllers) RegisterRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{AllowOrigins: "*", AllowHeaders: "Origin, Content-Type, Accept"}))
	app.Post("/register", c.register)
	app.Post("/login", c.login)
	app.Post("/stock", c.addStock)
	app.Get("/ws/stocks", websocket.New(c.GetAllStocks))
	app.Get("/ws/:ticker", websocket.New(c.GetStock))
	app.Post("/order", c.NewOrder)
}
