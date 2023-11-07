package user

import (
	"context"
)

type Repository interface {
	Set(ctx context.Context, user *User) error
	Get(ctx context.Context, uuid string) (*User, error)
	Update(ctx context.Context, user *User)
	Del(ctx context.Context, uuid string) error
}