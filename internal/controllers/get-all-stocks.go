package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

// TODO: remove magic sleep timeout
func (c *Controllers) GetAllStocks(ws *websocket.Conn) {
	const op = "controllers.GetAllStocks"
	defer func() {
		err := ws.Close()
		if err != nil {
			log.Printf("%s: %v\n", op, err)
		}
		log.Printf("WS conneciton closed")
	}()
	for {
		allStocksResponse, err := c.stockService.Stocks(c.ctx)
		if err != nil {
			log.Printf("%s: %v\n", op, err)
			ws.Close()
		}
		ws.WriteJSON(allStocksResponse)
		time.Sleep(1 * time.Second)
		fmt.Println(allStocksResponse)
	}
}
