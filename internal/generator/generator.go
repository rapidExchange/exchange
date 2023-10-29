package generator

import (
	"math"
	"math/rand"

	"rapidEx/internal/domain/stock"
	"rapidEx/internal/domain/order"
	"rapidEx/internal/stock-price-processor"
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
	stockProcessor := stockPriceProcessor.New()
	// prec := stockProcessor.PreciseAs(price)
	price = stockProcessor.Round(price, 10000)

	ord := order.Order{Ticker: s.Ticker, Quantity: volume, Price: price}

	switch ord.Price > s.Price {
	case true:
		s.Stockbook.Sell[ord.Price] += ord.Quantity
	default:
		s.Stockbook.Buy[ord.Price] += ord.Quantity
	}
	return ord
}