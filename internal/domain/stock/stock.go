package stock

import "rapidEx/internal/domain/stock-book"

type Stock struct{
	stockName string
	price	float64
	stockbook *stockBook.StockBook
}

func New(stockName string, price float64, stockBook *stockBook.StockBook) *Stock{
	return &Stock{
		stockName: stockName,
		price: price,
		stockbook: stockBook,
	}
} 