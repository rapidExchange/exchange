package stock

import "rapidEx/internal/domain/stock-book"

type Stock struct{
	Ticker string
	Price	float64
	Stockbook stockBook.StockBook
}

func New(ticker string, price float64) *Stock{
	return &Stock{
		Ticker: ticker,
		Price: price,
		Stockbook: *stockBook.New(),
	}
} 