package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"
)

//TODO: normal error handling

func RegisterRoutes(app *fiber.App) {
	app.Post("/register", register)
	app.Post("/login", login)
	app.Get("/all-tickers", GetAllTickers)
	app.Post("/add-stock", addStock)
	app.Get("/ws/get-all-stocks", websocket.New(GetAllStocks))
	app.Get("/ws/:ticker", websocket.New(GetStock))
	app.Post("/order", NewOrder)
}