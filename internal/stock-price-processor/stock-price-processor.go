package stockPriceProcessor

import (
	"context"
	"errors"
	"log"
	"math"
	"rapidEx/internal/domain/stock"
	stockBook "rapidEx/internal/domain/stock-book"
	redisconnect "rapidEx/internal/redis-connect"
	tickerstorage "rapidEx/internal/tickerStorage"
	"strings"
)

type stockPriceProcessor struct {
}

func New() *stockPriceProcessor {
	return &stockPriceProcessor{}
}

// Returns the meanweight value and its weight
func meanWeight(book map[float64]float64) (float64, float64) {
	valWeight, sumquant := 0.0, 0.0
	for val, quant := range book {
		valWeight += val * quant
		sumquant += quant
	}
	meanw := 0.0
	if sumquant != 0 {
		meanw = valWeight / sumquant
	}
	return meanw, sumquant
}

// Returns the meanweight price of stockbook. With an empty stockbook it returns zero and an error.
func (proc *stockPriceProcessor) MeanWeight(stockBook stockBook.StockBook) (float64, error) {
	buyMean, buyWeight := meanWeight(stockBook.Buy)
	sellMean, sellWeight := meanWeight(stockBook.Sell)

	meanw := 0.0
	var err error
	if buyWeight != 0 || sellWeight != 0 {
		meanw = (buyMean*buyWeight + sellMean*sellWeight) / (buyWeight + sellWeight)
		err = nil
	} else {
		err = errors.New("empty stock book")
	}
	return meanw, err
}

func (proc *stockPriceProcessor) Round(x float64, prec int) float64 {
	pow := math.Pow10(prec)
	rounded := math.Floor(x * pow)
	return rounded / pow
}

func (proc *stockPriceProcessor) PreciseAs(s string) int {
	st := strings.TrimRight(strings.Split(s, ".")[1], "0")
	return len(st)
}

// UpdatePrices updates all stocks price
func (proc *stockPriceProcessor) UpdatePrices() {
	tickerStorage := tickerstorage.GetInstanse()
	tickers := tickerStorage.GetTickers()
	for _, ticker := range tickers {
		err := proc.UpdatePrice(ticker)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (proc *stockPriceProcessor) UpdatePrice(ticker string) error {
	redisClient, err := redisconnect.SetRedisConn()
	if err != nil {
		return err
	}
	sRep := stock.NewRepository(redisClient)
	ctx := context.Background()

	stock, err := sRep.Get(ctx, ticker)
	if err != nil {
		return err
	}
	price, err := proc.MeanWeight(stock.Stockbook)
	if err != nil {
		return err
	}
	stock.Price = price

	log.Printf("New price of %s : %.20f", ticker, stock.Price)

	err = sRep.Set(ctx, *stock)
	return err
}
