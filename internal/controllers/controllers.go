package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func home(c *fiber.Ctx) error {
	return c.SendString("Start page")
}

func stocks(c *fiber.Ctx) error {
	return c.SendString("List of stocks page")
}

func stockPair(c *fiber.Ctx) error {
	return c.SendString("Pair of stocks trade page")
}

func login(c *fiber.Ctx) error {
	return c.SendString("Login route")
}

func register(c *fiber.Ctx) error {
	return c.SendString("Register route")
}

func RegisterRoutes() {
	app := fiber.New()

	app.Post("/register", register)
	app.Post("/login", login)
	app.Get("/stocks/:pair", stocks)
	app.Get("/stocks", stockPair)
	app.Get("/", home)
}