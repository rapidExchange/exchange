package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func stocks(c *fiber.Ctx) error {
	return c.SendString("List of stocks page")
}