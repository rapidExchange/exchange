package controllers

import (
	"encoding/json"
	"errors"
	tickerstorage "rapidEx/internal/tickerStorage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type AllTickersResponse struct {
	Tickers []string `json:"tickers"`
}

func GetAllTickers(c *fiber.Ctx) error {
	const op = "controllers.GetAllTickers"
	storage := tickerstorage.GetInstanse()
	tickers := storage.GetTickers()

	if !validateTickers(tickers) {
		return errors.New("no tickers")
	}
	tickersResponse := &AllTickersResponse{Tickers: tickers}
	bytes, err := json.Marshal(tickersResponse)
	if err != nil {
		log.Errorf("%s: %w", op, err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Send(bytes)
}

func validateTickers(tickers []string) bool {
	return len(tickers) == 0
}
