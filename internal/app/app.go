package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"rapidEx/config"
	"rapidEx/internal/controllers"
	stockDomain "rapidEx/internal/domain/stock"
	redisconnect "rapidEx/internal/redis"
	stockrepository "rapidEx/internal/repositories/stock-repository"
	"rapidEx/internal/services/stock"
	tickerstorage "rapidEx/internal/tickerStorage"
	"time"

	"github.com/gofiber/fiber/v2"
)

type app struct {
	c                   config.Config
	ctx                 context.Context
	gen                 Generator
	dealsProcessor      DealsProcessor
	stockPriceProcessor StockPriceProcessor
}

type Generator interface {
	GenerateALot(stock *stockDomain.Stock, genNum int)
}

type DealsProcessor interface {
	Do(s *stockDomain.Stock)
}

type StockPriceProcessor interface {
	UpdatePrice(stock *stockDomain.Stock) error
}

type StockProvider interface {
	Set(ctx context.Context, stock *stockDomain.Stock) error
	Stock(ctx context.Context, ticker string) (*stockDomain.Stock, error)
}

func New(gen Generator,
	dealsDealsProcessor DealsProcessor,
	stockPriceProcessor StockPriceProcessor) (*app, error) {
	pwd, _ := os.Getwd()
	c, err := config.LoadConfig(pwd)
	if err != nil {
		return nil, err
	}
	return &app{c: c, gen: gen, dealsProcessor: dealsDealsProcessor,
		stockPriceProcessor: stockPriceProcessor}, nil

}

func (a *app) Do() {
	for {
		tickers := getTickers()
		if len(tickers) == 0 {
			log.Println("Waiting stocks for generate")
		}
		stockRepository := stockrepository.NewStockRepository(redisconnect.MustConnect())
		stockProvider := stock.New(slog.New(slog.NewTextHandler(os.Stdout, nil)), stockRepository, stockRepository, nil)

		for _, ticker := range tickers {
			stock, err := stockProvider.Stock(context.Background(), ticker)
			if err != nil {
				log.Println(err)
				continue
			}
			go a.handleStock(stock, stockProvider)
			log.Printf("new price of %s: %f\n", stock.Ticker, stock.Price)
		}
		time.Sleep(1 * time.Second)
	}
}

func (a *app) handleStock(stock *stockDomain.Stock, stockProvider StockProvider) {
	a.gen.GenerateALot(stock, 10)
	a.dealsProcessor.Do(stock)
	a.stockPriceProcessor.UpdatePrice(stock)
	if err := stockProvider.Set(context.Background(), stock); err != nil {
		log.Println(err)
	}
}

func (a *app) ListenAndServe() {
	fiberApp := fiber.New()

	controllers.RegisterRoutes(fiberApp)

	log.Fatal(fiberApp.Listen(fmt.Sprintf(":%s", a.c.AppPort)))
}

func getTickers() []string {
	tickerStorage := tickerstorage.GetInstanse()
	return tickerStorage.GetTickers()
}
