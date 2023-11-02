package stock_usecases

import (
	"context"
	"fmt"
	"os"
	"rapidEx/internal/domain/stock"

	"rapidEx/config"

	"github.com/redis/go-redis/v9"
)

func SetStock(ticker string, price float64) error {
	s := stock.New(ticker, price)

	rc, err := setRedisConn()

	if err != nil {
		return err
	}

	stockRepo := stock.NewRepository(rc)

	ctx := context.Background()

	err = stockRepo.Set(ctx, *s)
	if err != nil {
		return err
	}
	return nil
}

func GetStock(ticker string) (*stock.Stock, error) {
	rc, err := setRedisConn()
	if err != nil {
		return nil, err
	}
	stockRepo := stock.NewRepository(rc)

	ctx := context.Background()

	stock, err := stockRepo.Get(ctx, ticker)
	if err != nil {
		return nil, err
	}
	return stock, nil
}
func setRedisConn() (*redis.Client, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err 
	}

	c, err := config.LoadConfig(pwd)

	if err != nil {
		return nil, err
	}

	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@localhost:6379/1", c.RedisUser, c.RedisPassword))
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opt), nil
}