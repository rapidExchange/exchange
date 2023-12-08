package user

import (
	"github.com/google/uuid"
)

type User struct {
	UUID           uuid.UUID
	Email          string
	PasswordHash   string
	Balance map[string]float64
}

func New(email, passwordhash string) *User {
	return &User{
		UUID:           uuid.New(),
		Email:          email,
		PasswordHash:   passwordhash,
		Balance: make(map[string]float64),
	}
}
