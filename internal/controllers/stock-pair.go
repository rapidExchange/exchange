package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func stockPair(c *fiber.Ctx) error {
	return c.SendString("Pair of stocks trade page")
}