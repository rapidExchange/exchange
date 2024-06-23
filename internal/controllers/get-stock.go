package controllers

import (
	"context"
	"log"
	"strconv"
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

// TODO: remove magic sleep timeout
func (c *Controllers) GetStock(ws *websocket.Conn) {
	const op = "controllers.GetStock"
	defer func() {
		err := ws.Close()
		if err != nil {
			log.Printf("%s: %v\n", op, err)
		}
		log.Println("WS conneciton closed")
	}()
	ticker := strings.ReplaceAll(ws.Params("ticker"), "_", "/")
	for {
		stock, err := c.stockService.Stock(context.Background(), ticker)
		if err != nil {
			log.Printf("%s: %v\n", op, err)
			ws.Close()
		}
		ws.WriteJSON(getStockWS{Ticker: stock.Ticker,
			Price: stock.Price, Buy: mapFloatToString(stock.Stockbook.Buy),
			Sell: mapFloatToString(stock.Stockbook.Sell)})
		time.Sleep(time.Second)
	}
}

func mapFloatToString(m map[float64]float64) map[string]string {
	stringMap := make(map[string]string)

	for k, v := range m {
		stringKey := strconv.FormatFloat(k, 'f', -1, 64)
		stringValue := strconv.FormatFloat(v, 'f', -1, 64)
		stringMap[stringKey] = stringValue
	}
	return stringMap
}
