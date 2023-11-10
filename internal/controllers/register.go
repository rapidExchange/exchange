package controllers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"rapidEx/internal/domain/user"
	mysqlconnect "rapidEx/internal/mysql-connect"
)

type registerRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func register(c *fiber.Ctx) error {
	registerReq := new(registerRequest)
	if err := c.BodyParser(&registerReq); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	mc, err := mysqlconnect.SetMysqlConnection()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	userRepository := user.NewRepository(mc)
	b, err := registerCheck(registerReq.Email, userRepository)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if b {
		return c.SendStatus(fiber.StatusConflict)
	}
	usr := user.New(registerReq.Email, registerReq.PasswordHash)
	err = userRepository.Create(context.Background(), usr)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)

}
