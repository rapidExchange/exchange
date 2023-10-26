package controllers

import (
	"io"
	"net/http"
	"strings"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

//TODO: normal error handling

type getTickerPriceRequest struct {
	Symbol string `json:"symbol"`
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

func panel(c *fiber.Ctx) error {
	return c.SendString("Admin panel")
}

func getTickerPrice(c *fiber.Ctx) error {
	getPrice := new(getTickerPriceRequest)

	if err := c.BodyParser(getPrice); err != nil {
		return err
	}

	url := "https://api.binance.com/api/v3/ticker/price?symbol=" + strings.ToUpper(strings.TrimSpace(getPrice.Symbol))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	priceResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var binanceResp getTickerPriceBinanceResponse

	err = json.Unmarshal(priceResp, &binanceResp)
	if err != nil {
		return err
	}

	return c.SendString(binanceResp.Price)
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
	app.Get("/panel", panel)
	app.Post("/price", getTickerPrice)
	app.Get("/", home)
}