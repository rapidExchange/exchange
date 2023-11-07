package order

import "encoding/json"

type Order struct {
	Ticker   string
	Quantity float64
	Price    float64
	User     string
}

func New(ticker, user string, quantity, price float64) *Order {
	return &Order{Ticker: ticker,
		User:     user,
		Quantity: quantity,
		Price:    price}
}

func (o *Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Order) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}
