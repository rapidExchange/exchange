package stock

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Set(ctx context.Context, stock Stock) error
	Get(ctx context.Context, ticker string) (*Stock, error)
	Del(ctx context.Context, ticker string) error
}

//Redis client
type rsClient struct {
	rc	*redis.Client
}

func (r *rsClient) Set(ctx context.Context, stock Stock) error {
	status := r.rc.Set(ctx, stock.Ticker, stock, 1)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *rsClient) Get(ctx context.Context, ticker string) (*Stock, error) {
	rStock := r.rc.Get(ctx, ticker)
	if rStock.Err() != nil {
		return nil, rStock.Err()
	}

	s := Stock{}

	err := rStock.Scan(&s)
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