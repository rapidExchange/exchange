package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func register(c *fiber.Ctx) error {
	return c.SendString("Register route")
}