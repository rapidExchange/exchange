package order

type Order struct {
	ticker	string
	quantity	float64
	price	float64
}

func New(ticker string, quantity, price float64) Order {
	return Order{ticker: ticker,
		quantity: quantity,
		price: price,}
}