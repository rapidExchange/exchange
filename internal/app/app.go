package app

import (
	"context"
	"log"
	"os"
	"rapidEx/config"
	"rapidEx/internal/controllers"
	"rapidEx/internal/domain/stock"

	"github.com/gofiber/fiber/v2"
)

type app struct {
	c                   config.Config
	ctx                 context.Context
	gen                 Generator
	dealsProcessor      DealsProcessor
	stockPriceProcessor StockPriceProcessor
	tickerStorage       TickerStorage
}

type Generator interface {
	GenerateALot(stock *stock.Stock, genNum int)
}

type DealsProcessor interface {
	Do()
}

type StockPriceProcessor interface {
	UpdatePrice(stock *stock.Stock) error
}

type TickerStorage interface {
	GetTickers() []string
}

func New(gen Generator,
	dealsDealsProcessor DealsProcessor,
	stockPriceProcessor StockPriceProcessor,
	tickerStorage TickerStorage) (*app, error) {
	pwd, _ := os.Getwd()
	c, err := config.LoadConfig(pwd)
	if err != nil {
		return nil, err
	}

	return &app{c: c, gen: gen, dealsProcessor: dealsDealsProcessor,
		stockPriceProcessor: stockPriceProcessor, tickerStorage: tickerStorage}, nil

}

func (a *app) Do() {

}

func (a *app) ListenAndServe() {
	fiberApp := fiber.New()

	controllers.RegisterRoutes(fiberApp)

	log.Fatal(fiberApp.Listen(":2345"))
}
