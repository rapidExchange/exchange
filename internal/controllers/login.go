package controllers

import (
	"context"
	"errors"
	"log"

	"rapidEx/internal/storage"

	"github.com/gofiber/fiber/v2"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Controllers) login(ctx *fiber.Ctx) error {
	loginReq := loginRequest{}
	if err := ctx.BodyParser(&loginReq); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	token, err := c.authService.Login(context.Background(), loginReq.Email, loginReq.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendString(token)
}
