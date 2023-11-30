package user

import "github.com/google/uuid"

type User struct {
	UUID           uuid.UUID
	Email          string
	PasswordHash   string
	OrdersQuantity int
	Balance Balance
}

// bad solution for store balance
type Balance struct {
	balanceSheet map[string]float64
}

func New(email, passwordhash string) *User {
	return &User{
		UUID:           uuid.New(),
		Email:          email,
		PasswordHash:   passwordhash,
		OrdersQuantity: 0,
		Balance: Balance{},
	}
}
