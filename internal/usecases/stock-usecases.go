package stockUsecases

import (
	"context"
	"rapidEx/internal/domain/stock"

	"github.com/redis/go-redis/v9"
)

//TODO: remove dependence from stock repository
func AddStock(rc *redis.Client, ticker string, price float64) error {
	s := stock.New(ticker, price)

	stockRepo := stock.NewRepository(rc)

	var ctx context.Context
	err := stockRepo.Set(ctx, *s)
	if err != nil {
		return err
	}
	return nil
}