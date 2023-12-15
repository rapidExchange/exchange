package generator

import (
	"context"
	"log"
	"math"
	"math/rand"

	"rapidEx/internal/domain/order"
	"rapidEx/internal/domain/stock"
	"rapidEx/internal/redis-connect"
	stockrepository "rapidEx/internal/repositories/stock-repository"
	"rapidEx/internal/tickerStorage"
	"rapidEx/internal/utils"
)

type Generator struct {
}

func New() *Generator {
	return &Generator{}
}

func (g Generator) generate(cPrice float64, ticker string) (float64, float64) {
	volume := float64(rand.Int31n(5000-10)+10) + rand.Float64()
	price := math.Abs(rand.NormFloat64()*0.05*cPrice + cPrice)
	prec, ok := tickerstorage.GetInstanse().Get(ticker)
	if !ok {
		log.Println("generator: undefined ticker: ", ticker)
		return 0, 0
	}
	return utils.Round(volume, prec+3), utils.Round(price, prec)
}

func (g Generator) OrderGenerate(s *stock.Stock) order.Order {
	volume, price := g.generate(s.Price, s.Ticker)

	Order := order.Order{Ticker: s.Ticker, Quantity: volume, Price: price}

	switch Order.Price > s.Price {
	case true:
		s.Stockbook.Sell[Order.Price] += Order.Quantity
	default:
		s.Stockbook.Buy[Order.Price] += Order.Quantity
	}
	return Order
}

func (g Generator) GenerateForAll() {
	tickerStorage := tickerstorage.GetInstanse()
	tickers := tickerStorage.GetTickers()
	if len(tickers) == 0 {
		log.Println("Waiting stocks for orders generate")
	}
	for _, ticker := range tickers {
		redisClient := redisconnect.MustConnect()

		stockRepository := stockrepository.NewStockRepository(redisClient)

		ctx := context.Background()

		stock, err := stockRepository.Get(ctx, ticker)
		if err != nil {
			log.Println(err)
			continue
		}

		g.GenerateALot(stock, 10)

		stockRepository.Set(ctx, *stock)
		
	}
}

//GenerateALot generates a lot orders for one stock sequentially
func (g Generator)GenerateALot(s *stock.Stock, genNum int) {
	for i := 0; i < genNum; i++ {
		g.OrderGenerate(s)
	}
}