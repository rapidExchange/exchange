package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"rapidEx/internal/domain/user"
	mysqlconnect "rapidEx/internal/mysql-connect"
	userrepository "rapidEx/internal/repositories/user-repository"
)

type registerRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func register(c *fiber.Ctx) error {
	const op = "controllers.register"
	registerReq := new(registerRequest)
	if err := c.BodyParser(&registerReq); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	mc := mysqlconnect.MustConnect()
	userRepository := userrepository.NewUserRepository(mc)
	b, err := registerCheck(registerReq.Email, userRepository)
	if err != nil {
		log.Println(fmt.Errorf("%s: %w", op, err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if b {
		return c.SendStatus(fiber.StatusConflict)
	}
	usr := user.New(registerReq.Email, registerReq.PasswordHash)
	err = userRepository.Create(context.Background(), usr)
	if err != nil {
		log.Println(fmt.Errorf("%s: %w", op, err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)

}
