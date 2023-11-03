package stock_usecases

import (
	"context"
	"rapidEx/internal/domain/stock"

	"rapidEx/internal/utils"
)

func SetStock(ticker string, price float64) error {
	s := stock.New(ticker, price)

	rc, err := utils.SetRedisConn()

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
	rc, err := utils.SetRedisConn()
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
