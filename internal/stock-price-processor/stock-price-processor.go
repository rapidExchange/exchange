package stockPriceProcessor

import (
	"context"
	"errors"
	"fmt"
	"math"
	"rapidEx/internal/domain/stock"
	stockBook "rapidEx/internal/domain/stock-book"
	"rapidEx/internal/tickerStorage"
	"rapidEx/internal/usecases/stock_usecases"
	"rapidEx/internal/utils"
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
	fmt.Println(len(stockBook.Buy), len(stockBook.Sell))

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
			fmt.Println(err)
			continue
		}
	}
}

func (proc *stockPriceProcessor) UpdatePrice(ticker string) error {
	s, err := stock_usecases.GetStock(ticker)
	if err != nil {
		return err
	}
	price, err := proc.MeanWeight(s.Stockbook)
	if err != nil {
		return err
	}
	s.Price = price
	rc, err := utils.SetRedisConn()
	if err != nil {
		return err
	}
	sRep := stock.NewRepository(rc)
	
	ctx := context.Background()

	fmt.Printf("New price of %s : %f", ticker, s.Price)

	err = sRep.Set(ctx, *s)
	return err
}