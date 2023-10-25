package stock

import "rapidEx/internal/domain/stock-book"

type Stock struct{
	ticker string
	price	float64
	stockbook *stockBook.StockBook
}

func New(ticker string, price float64, stockBook *stockBook.StockBook) *Stock{
	return &Stock{
		ticker: ticker,
		price: price,
		stockbook: stockBook,
	}
} 