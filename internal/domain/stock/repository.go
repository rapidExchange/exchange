package stock

import (
	"context"
	"encoding/json"
	"time"
	"errors"

	"github.com/redis/go-redis/v9"

	"rapidEx/internal/utils"
)

type StockMapString struct {
	Ticker string
	Price float64
	Buy map[string]string
	Sell map[string]string
}

func NewStockMapString(s Stock) *StockMapString {
	var sMap StockMapString
	sMap.Ticker = s.Ticker
	sMap.Price = s.Price
	sMap.Buy = utils.MapFloatToString(s.Stockbook.Buy)
	sMap.Sell = utils.MapFloatToString(s.Stockbook.Sell)
	return &sMap
}

func (s StockMapString)MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}


func UnmarshalBinary(data []byte) (*StockMapString, error) {
	var s StockMapString
	err := json.Unmarshal(data, &s)
	return &s, err
}
//TODO: handle errors

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
	stockMapString := NewStockMapString(stock)
	status := r.rc.Set(ctx, stock.Ticker, stockMapString, time.Second * 1000)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (r *rsClient) Get(ctx context.Context, ticker string) (*Stock, error) {
	rStockMapString := r.rc.Get(ctx, ticker)
	switch {
	case rStockMapString.Err() == redis.Nil:
		return nil, errors.New("Stock not found")
	case rStockMapString.Err() != nil:
		return nil, rStockMapString.Err()
	}

	var s Stock

	sStockWrap, err := rStockMapString.Result()
	if err != nil {
		return nil, err
	}

	stockMapString, err := UnmarshalBinary([]byte(sStockWrap))
	if err != nil {
		return nil, err
	}

	mBuy, err:= utils.MapStringToFloat(stockMapString.Buy)
	if err != nil {
		return nil, err
	}

	mSell, err := utils.MapStringToFloat(stockMapString.Sell)
	if err != nil {
		return nil, err
	}

	s.Ticker = stockMapString.Ticker
	s.Price = stockMapString.Price
	s.Stockbook.Buy = mBuy
	s.Stockbook.Sell = mSell
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