package controllers

import (
	"context"
	"fmt"
	"log"
	"rapidEx/internal/domain/order"
	redisconnect "rapidEx/internal/redis-connect"
	orderrepository "rapidEx/internal/repositories/order-repository"

	"github.com/gofiber/fiber/v2"
)

type NewOrderRequest struct {
	Email    string  `json:"email"`
	Type     string  `json:"type"`
	Ticker   string  `json:"ticker"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

// TODO: check user balance
func NewOrder(c *fiber.Ctx) error {
	const op = "controllers.NewOrder"
	newOrderRequest := new(NewOrderRequest)
	if err := c.BodyParser(&newOrderRequest); err != nil {
		log.Println(fmt.Errorf("%s: %w", op, err))
		return c.SendStatus(fiber.StatusBadRequest)
	}
	order, err := order.New(newOrderRequest.Ticker, newOrderRequest.Email, newOrderRequest.Type, newOrderRequest.Quantity, newOrderRequest.Price)
	if err != nil {
		log.Println(fmt.Errorf("%s: %w", op, err))
		return c.SendStatus(fiber.StatusBadRequest)
	}
	redisClient, err := redisconnect.SetRedisConn()
	if err != nil {
		log.Println(fmt.Errorf("%s: %w", op, err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	orderRepository := orderrepository.NewOrderRepository(redisClient)
	err = orderRepository.Set(context.Background(), order)
	if err != nil {
		log.Println(fmt.Errorf("%s: %w", op, err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
