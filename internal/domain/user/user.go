package user

import "github.com/google/uuid"

type User struct {
	UUID           uuid.UUID
	Username       string
	Email          string
	PasswordHash   string
	OrdersQuantity int
}

func New(username, email, passwordhash string) *User {
	return &User{
		UUID:           uuid.New(),
		Username:       username,
		Email:          email,
		PasswordHash:   passwordhash,
		OrdersQuantity: 0,
	}
}
