package order

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Set(ctx context.Context, order Order) error
	Get(ctx context.Context, ticker string) (*Order, error)
	Del(ctx context.Context, ticker string) error
}

//Redis client
type rsClient struct {
	rc	*redis.Client
}

func (r *rsClient) Set(ctx context.Context, order Order) error {
	status := r.rc.Set(ctx, order.ticker, order, 1)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *rsClient) Get(ctx context.Context, ticker string) (*Order, error) {
	rOrder := r.rc.Get(ctx, ticker)
	if rOrder.Err() != nil {
		return nil, rOrder.Err()
	}

	s := Order{}

	err := rOrder.Scan(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *rsClient) Del(ctx context.Context, ticker string) error {
	intCmd := r.rc.Del(ctx, ticker)

	if intCmd.Err() != nil {
		return intCmd.Err()
	}

	return nil
}

func NewRepository(rc *redis.Client) Repository {
	return &rsClient{rc: rc}
}