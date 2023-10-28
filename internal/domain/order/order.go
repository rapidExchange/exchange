package order

type Order struct {
	Ticker	string
	Quantity	float64
	Price	float64
}

func New(ticker string, quantity, price float64) Order {
	return Order{Ticker: ticker,
		Quantity: quantity,
		Price: price,}
}