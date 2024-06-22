package controllers

import (
	"encoding/json"
	"errors"
	"log"
	tickerstorage "rapidEx/internal/tickerStorage"

	"github.com/gofiber/fiber/v2"
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
		log.Printf("%s: %w\n", op, err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Send(bytes)
}

func validateTickers(tickers []string) bool {
	return len(tickers) == 0
}
