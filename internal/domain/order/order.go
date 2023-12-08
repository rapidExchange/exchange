package order

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type Order struct {
	OrderUUID uuid.UUID
	Ticker    string
	Quantity  float64
	Price     float64
	Email     string
	Status    string
	Type      string // buy or sell(b or s) order type
}

func New(ticker, email, Type string, quantity, price float64) (*Order, error) {
	typeCheck := Type == "b" || Type == "s"
	if !typeCheck {
		return nil, errors.New("unsopported order type")
	}
	return &Order{
		OrderUUID: uuid.New(),
		Ticker:    ticker,
		Email:     email,
		Quantity:  quantity,
		Price:     price,
		Status: "new",
		Type:      Type}, nil
}

func (o Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Order) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}
