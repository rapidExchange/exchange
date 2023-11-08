package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"rapidEx/internal/domain/stock"
	redisconnect "rapidEx/internal/redis-connect"
	stockPriceProcessor "rapidEx/internal/stock-price-processor"
	tickerstorage "rapidEx/internal/tickerStorage"
)

type getTickerPriceBinanceRequest struct {
	FirstSymbol  string `json:"first_symbol"`
	SecondSymbol string `json:"second_symbol"`
}

type getTickerPriceBinanceResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

func addStock(c *fiber.Ctx) error {
	getPriceBinanceRequest := new(getTickerPriceBinanceRequest)
	if err := c.BodyParser(&getPriceBinanceRequest); err != nil {
		return err
	}
	symbol := createSymbol(getPriceBinanceRequest.FirstSymbol, getPriceBinanceRequest.SecondSymbol)
	price, err := getBinancePrice(symbol)
	if err != nil {
		return err
	}
	precision := getPrecision(price)
	roundedPrice := roundWithPrecision(price)
	ticker := createTicker(getPriceBinanceRequest.FirstSymbol, getPriceBinanceRequest.SecondSymbol)
	err = setStock(ticker, roundedPrice)
	if err != nil {
		return err
	}
	setTickerToStorage(ticker, precision)
	return c.SendStatus(fiber.StatusOK)
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

func setStock(ticker string, price float64) error {
	Stock := stock.New(ticker, price)
	redisClient, err := redisconnect.SetRedisConn()
	if err != nil {
		return err
	}
	stockRepository := stock.NewRepository(redisClient)
	err = stockRepository.Set(context.Background(), *Stock)

	return err
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

func roundWithPrecision(price float64) float64 {
	priceProcessor := stockPriceProcessor.New()
	precision := getPrecision(price)
	return priceProcessor.Round(price, precision)
}

func getPrecision(price float64) int {
	priceProcessor := stockPriceProcessor.New()
	return priceProcessor.PreciseAs(strconv.FormatFloat(price, 'f', -1, 64))
}

func setTickerToStorage(ticker string, precision int) {
	tickerStorage := tickerstorage.GetInstanse()
	tickerStorage.TickerAppend(ticker, precision)
}

func readBody(source io.Reader) ([]byte, error) {
	return io.ReadAll(source)
}

func createSymbol(firstStock, secondStock string) string {
	return strings.ToUpper(firstStock + secondStock)
}

func createTicker(firstStock, secondStock string) string {
	return strings.ToLower(firstStock + "/" + secondStock)
}