package controllers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Controllers) register(ctx *fiber.Ctx) error {
	registerReq := new(registerRequest)
	if err := ctx.BodyParser(&registerReq); err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	err := c.authService.Register(context.Background(), registerReq.Email, registerReq.Password)
	if err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}
