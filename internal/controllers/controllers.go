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

//TODO: normal error handling

type getTickerPriceBinanceRequest struct {
	FirstSymbol  string `json:"first_symbol"`
	SecondSymbol string `json:"second_symbol"`
	Precision    int    `json:"precision"`
}

type getTickerPriceBinanceResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
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
	getPriceBinanceRequest := new(getTickerPriceBinanceRequest)

	if err := c.BodyParser(&getPriceBinanceRequest); err != nil {
		return err
	}

	symbol := createSymbol(getPriceBinanceRequest.FirstSymbol, getPriceBinanceRequest.SecondSymbol)

	price, err := getBinancePrice(symbol)
	if err != nil {
		return err
	}

	ticker := createTicker(getPriceBinanceRequest.FirstSymbol, getPriceBinanceRequest.SecondSymbol)

	err = setStock(ticker, price)
	if err != nil {
		return err
	}

	setTickerToStorage(ticker, getPriceBinanceRequest.Precision)

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
	proc := stockPriceProcessor.New()
	prec := proc.PreciseAs(strconv.FormatFloat(binanceResponse.Price, 'f', -1, 64)) // need to write in getPriceBinanceRequest.Precision (see addStock)
	binancePrice := proc.Round(binanceResponse.Price, prec)
	return binancePrice, nil
}

func setStock(ticker string, price float64) error {
	Stock := stock.New(ticker, price)

	redisClient, err := redisconnect.SetRedisConn()
	if err != nil {
		return err
	}

	stockRepository := stock.NewRepository(redisClient)

	ctx := context.Background()

	err = stockRepository.Set(ctx, *Stock)

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
