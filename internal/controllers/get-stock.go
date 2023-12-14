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

	redisConneciton := redisconnect.MustConnect()
	stockRepository := stockrepository.NewStockRepository(redisConneciton)
	for {
	s, err := stockRepository.Get(context.Background(), ticker)
	if err != nil {
		log.Println(fmt.Errorf("%s: %w", op, err))
		continue
	}
	stockModify := stockrepository.NewStockMapString(*s)
	c.WriteJSON(stockModify)
	time.Sleep(time.Second)
}
}

func validateStock(stock string) bool {
	return stock != ""
}
