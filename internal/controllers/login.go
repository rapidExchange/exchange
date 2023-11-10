package controllers

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"rapidEx/internal/domain/user"
	"rapidEx/internal/mysql-connect"
)

type loginRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func login(c *fiber.Ctx) error {
	loginReq := loginRequest{}
	if err := c.BodyParser(&loginReq); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	mc, err := mysqlconnect.SetMysqlConnection()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	userRepository := user.NewRepository(mc)
	ifReg, err := registerCheck(loginReq.Email, userRepository)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !ifReg {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.SendStatus(fiber.StatusOK)
}

func registerCheck(email string, userRepository user.Repository) (bool, error) {
	_, err := userRepository.Get(context.Background(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
