package generator

import (
	"context"
	"log"
	"math"
	"math/rand"

	"rapidEx/internal/domain/order"
	"rapidEx/internal/domain/stock"
	"rapidEx/internal/tickerStorage"
	"rapidEx/internal/redis-connect"
)

type Generator struct {
}

func New() *Generator {
	return &Generator{}
}

func (g Generator) generate(cPrice float64) (volume float64, price float64) {
	volume = float64(rand.Int31n(5000-10)+10) + rand.Float64()
	price = math.Abs(rand.NormFloat64()*0.05*cPrice + cPrice) // editable compression parameter
	return
}

func (g Generator) OrderGenerate(s *stock.Stock) order.Order {
	volume, price := g.generate(s.Price)

	ord := order.Order{Ticker: s.Ticker, Quantity: volume, Price: price}

	switch ord.Price > s.Price {
	case true:
		s.Stockbook.Sell[ord.Price] += ord.Quantity
	default:
		s.Stockbook.Buy[ord.Price] += ord.Quantity
	}
	return ord
}

func (g Generator) GenerateForAll() {
	tickerStorage := tickerstorage.GetInstanse()
	tickers := tickerStorage.GetTickers()
	if len(tickers) == 0 {
		log.Println("Waiting stocks for orders generate")
	}
	for _, ticker := range tickers {
		redisClient, err := redisconnect.SetRedisConn()
		if err != nil {
			log.Println(err)
			continue
		}

		stockRepository := stock.NewRepository(redisClient)

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