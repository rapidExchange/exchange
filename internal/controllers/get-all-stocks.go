package controllers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	redisconnect "rapidEx/internal/redis"
	stockrepository "rapidEx/internal/repositories/stock-repository"
	"rapidEx/internal/services/stock"
	tickerstorage "rapidEx/internal/tickerStorage"
	"time"

	"github.com/gofiber/contrib/websocket"
)

func GetAllStocks(c *websocket.Conn) {
	const op = "controllers.GetAllStocks"
	defer func() {
		err := c.Close()
		if err != nil {
			log.Printf("%s: %w\n", op, err)
		}
		log.Printf("WS conneciton closed")
	}()

	tickers := tickerstorage.GetInstanse().GetTickers()

	redisConneciton := redisconnect.MustConnect()
	stockRepository := stockrepository.NewStockRepository(redisConneciton)
	stockMonitor := stock.New(slog.New(slog.NewTextHandler(os.Stdout, nil)), stockRepository,
		stockRepository, nil)
	for {
		allStocksResponse := make(map[string]float64)
		for _, ticker := range tickers {
			stock, err := stockMonitor.Stock(context.Background(), ticker)
			if err != nil {
				log.Println(fmt.Errorf("%s: %w", op, err))
				continue
			}
			allStocksResponse[stock.Ticker] = stock.Price
		}
		c.WriteJSON(allStocksResponse)
		time.Sleep(1 * time.Second)
		fmt.Println(allStocksResponse)
	}
}
