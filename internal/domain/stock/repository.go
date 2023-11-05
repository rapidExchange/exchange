package stock

import (
	"context"
	"encoding/json"
	"time"
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"

)

type StockModify struct {
	Ticker string
	Price float64
	Buy map[string]string
	Sell map[string]string
}

func NewStockMapString(s Stock) *StockModify {
	var sMap StockModify
	sMap.Ticker = s.Ticker
	sMap.Price = s.Price
	sMap.Buy = mapFloatToString(s.Stockbook.Buy)
	sMap.Sell = mapFloatToString(s.Stockbook.Sell)
	return &sMap
}

func (s StockModify)MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}


func UnmarshalBinary(data []byte) (*StockModify, error) {
	var s StockModify
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

	mBuy, err:=mapStringToFloat(stockMapString.Buy)
	if err != nil {
		return nil, err
	}

	mSell, err := mapStringToFloat(stockMapString.Sell)
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

func mapFloatToString(m map[float64]float64) map[string]string {
	stringMap := make(map[string]string)

	for k, v := range m {
		stringKey := strconv.FormatFloat(k, 'f', -1, 64)
		stringValue := strconv.FormatFloat(v, 'f', -1, 64)
		stringMap[stringKey] = stringValue
	}
	return stringMap
}

func mapStringToFloat(m map[string]string) (map[float64]float64, error) {
	floatMap := make(map[float64]float64)

	for k, v := range m {
		floatKey, err := strconv.ParseFloat(k, 64)
		if err != nil {
			return nil, err
		}
		floatVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}

		floatMap[floatKey] = floatVal
	}

	return floatMap, nil
}