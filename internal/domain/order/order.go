package order

type Order struct {
	stockName	string
	quantity	float64
	price	float64
}

func New(stockName string, quantity, price float64) Order {
	return Order{stockName: stockName,
		quantity: quantity,
		price: price,}
}