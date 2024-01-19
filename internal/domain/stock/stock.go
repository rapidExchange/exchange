package stock

import (
	stockBook "rapidEx/internal/domain/stock-book"
	"sync"
)

type Stock struct {
	Ticker    string
	Price     float64
	Stockbook stockBook.StockBook
	Lock sync.RWMutex
}

func New(ticker string, price float64) *Stock {
	return &Stock{
		Ticker:    ticker,
		Price:     price,
		Stockbook: *stockBook.New(),
		Lock: sync.RWMutex{},
	}
}
