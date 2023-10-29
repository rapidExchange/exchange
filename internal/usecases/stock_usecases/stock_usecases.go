package stock_usecases

import (
	"context"
	"fmt"
	"rapidEx/config"
	"rapidEx/internal/domain/stock"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

//TODO: remove dependence from stock repository
func SetStock(ticker string, price float64) error {
	s := stock.New(ticker, price)

	rc, err := setRedisConn()

	if err != nil {
		return err
	}

	stockRepo := stock.NewRepository(rc)

	var ctx context.Context
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

	var ctx context.Context
	stock, err := stockRepo.Get(ctx, ticker)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

func setRedisConn() (*redis.Client, error) {
	var c config.Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@localhost:6379", c.RedisUser, c.RedisPassword))

	if err != nil {
		return nil, err
	}

	return redis.NewClient(opt), nil
}