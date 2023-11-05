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
	tickerstorage "rapidEx/internal/tickerStorage"
)

//TODO: normal error handling

type getTickerPriceBinanceRequest struct {
	FirstSymbol  string `json:"first_symbol"`
	SecondSymbol string `json:"second_symbol"`
}

type getTickerPriceBinanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func home(c *fiber.Ctx) error {
	return c.SendString("Start page")
}

func stocks(c *fiber.Ctx) error {
	return c.SendString("List of stocks page")
}

func stockPair(c *fiber.Ctx) error {
	return c.SendString("Pair of stocks trade page")
}

func addStock(c *fiber.Ctx) error {
	getPriceRequest := new(getTickerPriceBinanceRequest)

	if err := c.BodyParser(&getPriceRequest); err != nil {
		return err
	}

	symbol := strings.ToUpper(strings.TrimSpace(getPriceRequest.FirstSymbol) +
		strings.TrimSpace(getPriceRequest.SecondSymbol))

	priceString, err := getBinancePrice(symbol)
	if err != nil {
		return err
	}

	price, err := strconv.ParseFloat(priceString, 64)
	if err != nil {
		return err
	}

	ticker := getPriceRequest.FirstSymbol + "/" + getPriceRequest.SecondSymbol

	Stock := stock.New(ticker, price)

	redisClient, err := redisconnect.SetRedisConn()
	if err != nil {
		return err
	}

	stockRepository := stock.NewRepository(redisClient)

	ctx := context.Background()

	stockRepository.Set(ctx, *Stock)

	tickerStorage := tickerstorage.GetInstanse()
	tickerStorage.TickerAppend(ticker)

	return c.SendStatus(fiber.StatusOK)
}

func login(c *fiber.Ctx) error {
	return c.SendString("Login route")
}

func register(c *fiber.Ctx) error {
	return c.SendString("Register route")
}

func RegisterRoutes(app *fiber.App) {
	app.Post("/register", register)
	app.Post("/login", login)
	app.Get("/stocks/:pair", stocks)
	app.Get("/stocks", stockPair)
	app.Post("/add-stock", addStock)
	app.Get("/", home)
}

func getBinancePrice(symbol string) (string, error) {
	url := "https://api.binance.com/api/v3/ticker/price?symbol=" + symbol

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	priceResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var binanceResponse getTickerPriceBinanceResponse

	err = json.Unmarshal(priceResponse, &binanceResponse)
	if err != nil {
		return "", err
	}
	return binanceResponse.Price, nil
}
