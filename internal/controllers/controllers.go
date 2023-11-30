package controllers

import (
	"github.com/gofiber/fiber/v2"
)

//TODO: normal error handling

func RegisterRoutes(app *fiber.App) {
	app.Post("/register", register)
	app.Post("/login", login)
	app.Get("/all-tickers", GetAllTickers)
	app.Post("/add-stock", addStock)
	app.Get("")
}