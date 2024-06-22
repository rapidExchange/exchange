package stockPriceProcessor

import (
	"errors"
	"rapidEx/internal/domain/stock"
	stockBook "rapidEx/internal/domain/stock-book"
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

func (proc *stockPriceProcessor) UpdatePrice(stock *stock.Stock) error {
	price, err := proc.MeanWeight(stock.Stockbook)
	if err != nil {
		return err
	}
	stock.Price = price
	return nil
}
