package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"rapidEx/internal/test_case"
	"rapidEx/internal/usecases/stock_usecases"
)

//TODO: normal error handling

type getTickerPriceBinanceRequest struct {
	FirstSymbol string `json:"first_symbol"`
	SecondSymbol string `json:"second_symbol"`
}

type getTickerPriceBinanceResponse struct {
	Symbol string `json:"symbol"`
	Price string `json:"price"`
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
	go case1.Case()
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

	err = stock_usecases.SetStock(ticker, price)
	if err != nil {
		return err
	}

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
	app.Post("/price", addStock)
	app.Get("/", home)
}

func getBinancePrice(symbol string) (string, error) {
	url := "https://api.binance.com/api/v3/ticker/price?symbol=" + symbol

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	priceResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var binanceResp getTickerPriceBinanceResponse

	err = json.Unmarshal(priceResp, &binanceResp)
	if err != nil {
		return "", err
	}
	return binanceResp.Price, nil
} 