package controllers

import (
	"context"
	"fmt"
	"log"
	redisconnect "rapidEx/internal/redis-connect"
	stockrepository "rapidEx/internal/repositories/stock-repository"
	tickerstorage "rapidEx/internal/tickerStorage"

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
	ticker := c.Params("stock")

	if !validateStock(ticker) {
		log.Println("Empty stockname")
		return
	}
	storage := tickerstorage.GetInstanse()
	if !storage.Find(ticker) {
		log.Printf("Undefined ticker: %s", ticker)
	}

	redisConneciton := redisconnect.MustConnect()
	stockRepository := stockrepository.NewStockRepository(redisConneciton)
	s, err := stockRepository.Get(context.Background(), ticker)
	if err != nil {
		log.Println(fmt.Errorf("%s: %w", op, err))
		return
	}
	stockModify := stockrepository.NewStockMapString(*s)
	c.WriteJSON(stockModify)
}

func validateStock(stock string) bool {
	return stock == ""
}
