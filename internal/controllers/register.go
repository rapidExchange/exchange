package controllers

import (
	"context"
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	mysql "rapidEx/internal/mysql"
	userrepository "rapidEx/internal/repositories/user-repository"
	"rapidEx/internal/services/auth"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func register(c *fiber.Ctx) error {
	const op = "controllers.register"
	registerReq := new(registerRequest)
	if err := c.BodyParser(&registerReq); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	mc := mysql.MustConnect()
	userRepository := userrepository.NewUserRepository(mc)
	auth := auth.New(&slog.Logger{}, userRepository, userRepository)
	err := auth.Register(context.Background(), registerReq.Email, registerReq.Password)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
