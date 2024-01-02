package controllers

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	mysqlconnect "rapidEx/internal/mysql"
	userrepository "rapidEx/internal/repositories/user-repository"
	"rapidEx/internal/services/auth"
	"rapidEx/internal/storage"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(c *fiber.Ctx) error {
	loginReq := loginRequest{}
	if err := c.BodyParser(&loginReq); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	mc := mysqlconnect.MustConnect()
	userRepository := userrepository.NewUserRepository(mc)
	auth := auth.New(&slog.Logger{}, userRepository, userRepository)
	token, err := auth.Login(context.Background(), loginReq.Email, loginReq.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString(token)
}
