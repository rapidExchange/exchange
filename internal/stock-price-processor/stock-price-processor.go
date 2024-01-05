package stockPriceProcessor

import (
	"errors"
	"log"
	"rapidEx/internal/domain/stock"
	stockBook "rapidEx/internal/domain/stock-book"
	tickerstorage "rapidEx/internal/tickerStorage"
	"rapidEx/internal/utils"
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

func (proc *stockPriceProcessor) PreciseAs(s string) int {
	st := strings.TrimRight(strings.Split(s, ".")[1], "0")
	return len(st)
}

func (proc *stockPriceProcessor) UpdatePrice(stock *stock.Stock) error {
	price, err := proc.MeanWeight(stock.Stockbook)
	if err != nil {
		return err
	}
	tickerStorage := tickerstorage.GetInstanse()
	prec, ok := tickerStorage.Get(stock.Ticker)
	if !ok {
		log.Println("Undefined ticker in tickerstorage: ", stock.Ticker)
	}
	stock.Price = utils.Round(price, prec)

	log.Printf("New price of %s : %f", stock.Ticker, stock.Price)

	return nil
}
