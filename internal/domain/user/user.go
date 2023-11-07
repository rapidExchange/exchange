package user

import "github.com/google/uuid"

type User struct {
	UUID           uuid.UUID
	Email          string
	PasswordHash   string
	OrdersQuantity int
}

func New(username, email, passwordhash string) *User {
	return &User{
		UUID:           uuid.New(),
		Email:          email,
		PasswordHash:   passwordhash,
		OrdersQuantity: 0,
	}
}
