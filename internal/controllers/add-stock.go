package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"rapidEx/internal/domain/stock"
)

type getTickerPriceBinanceRequest struct {
	FirstSymbol  string `json:"first_symbol"`
	SecondSymbol string `json:"second_symbol"`
}

type getTickerPriceBinanceResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

func (c *Controllers)addStock(ctx *fiber.Ctx) error {
	const op = "controllers.addStock"
	getPriceBinanceRequest := new(getTickerPriceBinanceRequest)
	if err := ctx.BodyParser(&getPriceBinanceRequest); err != nil {
		log.Printf("%s: %v\n", op, err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	symbol := strings.ToUpper(getPriceBinanceRequest.FirstSymbol + getPriceBinanceRequest.SecondSymbol)
	price, err := getBinancePrice(symbol)
	if err != nil {
		log.Printf("%s: %v\n", op, err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	ticker := strings.ToLower(getPriceBinanceRequest.FirstSymbol + "/" + getPriceBinanceRequest.SecondSymbol)
	//provide context
	err = c.stockService.Set(context.Background(), stock.New(ticker, price))
	if err != nil {
		log.Printf("%s: %v\n", op, err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func getBinancePrice(symbol string) (float64, error) {
	url := "https://api.binance.com/api/v3/ticker/price?symbol=" + symbol
	var zero float64
	response, err := makeBinancePriceRequest(url)
	if err != nil {
		return zero, nil
	}
	priceBinanceResponse, err := readBody(response.Body)
	if err != nil {
		return zero, err
	}
	binanceResponse, err := unmarshalToBinanceResponse(priceBinanceResponse)
	if err != nil {
		return 0.0, err
	}
	return binanceResponse.Price, nil
}


func makeBinancePriceRequest(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func unmarshalToBinanceResponse(response []byte) (*getTickerPriceBinanceResponse, error) {
	binanceResponse := &getTickerPriceBinanceResponse{}
	if err := json.Unmarshal(response, &binanceResponse); err != nil {
		return nil, err
	}
	return binanceResponse, nil
}

func readBody(source io.Reader) ([]byte, error) {
	return io.ReadAll(source)
}