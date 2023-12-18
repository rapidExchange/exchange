package controllers

import (
	"context"
	"fmt"
	"log"
	redisconnect "rapidEx/internal/redis-connect"
	stockrepository "rapidEx/internal/repositories/stock-repository"
	tickerstorage "rapidEx/internal/tickerStorage"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type getStockWS struct {
	Ticker    string            `json:"ticker"`
	Price     float64           `json:"price"`
	Buy       map[string]string `json:"stockBookBuy"`
	Sell      map[string]string `json:"stockBookSell"`
	Precision int               `json:"precision"`
}

func GetStock(c *websocket.Conn) {
	const op = "controllers.GetStock"
	defer func() {
		err := c.Close()
		if err != nil {
			log.Println(fmt.Errorf("%s: %w", op, err))
		}
		log.Printf("WS conneciton closed")
	}()
	ticker := strings.ReplaceAll(c.Params("ticker"), "_", "/")

	if !validateStock(ticker) {
		log.Println("Empty stockname")
		return
	}
	storage := tickerstorage.GetInstanse()
	if !storage.Find(ticker) {
		log.Printf("Undefined ticker: %s", ticker)
		return
	}

	prec, _ := storage.Get(ticker)

	redisConneciton := redisconnect.MustConnect()
	stockRepository := stockrepository.NewStockRepository(redisConneciton)
	for {
		s, err := stockRepository.Get(context.Background(), ticker)
		if err != nil {
			log.Println(fmt.Errorf("%s: %w", op, err))
			continue
		}
		stockModify := stockrepository.NewStockMapString(*s)
		c.WriteJSON(getStockWS{Ticker: stockModify.Ticker,
			Price: stockModify.Price, Buy: stockModify.Buy,
			Sell: stockModify.Sell, Precision: prec})
		time.Sleep(time.Second)
	}
}

func validateStock(stock string) bool {
	return stock != ""
}
