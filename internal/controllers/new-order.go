package controllers

import (
	"context"
	"log"
	"rapidEx/internal/domain/order"

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
func (c *Controllers) NewOrder(ctx *fiber.Ctx) error {
	const op = "controllers.NewOrder"
	newOrderRequest := new(NewOrderRequest)
	if err := ctx.BodyParser(&newOrderRequest); err != nil {
		log.Printf("%s: %v\n", op, err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	order, err := order.New(newOrderRequest.Ticker, newOrderRequest.Email, newOrderRequest.Type, newOrderRequest.Quantity, newOrderRequest.Price)
	if err != nil {
		log.Printf("%s: %v\n", op, err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	order.Status = "processing"
	err = c.orderService.Set(context.Background(), order)
	if err != nil {
		log.Printf("%s: %v\n", op, err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}
