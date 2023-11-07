package user

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Set(ctx context.Context, user *User) error
	Get(ctx context.Context, uuid string) (*User, error)
	Del(ctx context.Context, uuid string) error
}

type rsClient struct {
	rc *redis.Client
}

func (r *rsClient) Set(ctx context.Context, user *User) error {
	return nil
}

func (r *rsClient) Get(ctx context.Context, uuid string) (*User, error) {
	return nil, nil
}

func (r *rsClient) Del(ctx context.Context, uuid string) error {
	return nil
}

func NewRepository(rc *redis.Client) Repository {
	return &rsClient{
		rc: rc,
	}
}
