package stockPriceProcessor

import (
	"context"
	"errors"
	"log"
	"math"
	"rapidEx/internal/domain/stock"
	stockBook "rapidEx/internal/domain/stock-book"
	"rapidEx/internal/tickerStorage"
	"rapidEx/internal/redis-connect"
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

// Rounds to a floor with given precision
func (proc *stockPriceProcessor) Round(x float64, prec int) float64 {
	pow := math.Pow10(prec)
	rounded := math.Floor(x * pow)
	return rounded / pow
}

// Returns the number of significant decimal places
func (proc *stockPriceProcessor) PreciseAs(part float64) int {
	k := 0
	if part < 0.99 {
		for int(part)%10 == 0 {
			k++
			part *= 10
		}
	}
	return k
}

//UpdatePrices updates all stocks price
func (proc *stockPriceProcessor) UpdatePrices()  {
	tickerStorage := tickerstorage.GetInstanse()
	tickers := tickerStorage.GetTickers()
	for _, ticker := range tickers {
		err := proc.UpdatePrice(ticker)
		if err != nil{
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
	
	log.Printf("New price of %s : %f", ticker, stock.Price)

	err = sRep.Set(ctx, *stock)
	return err
}