package stockrepository

import (
	"context"
	"encoding/json"
	"errors"
	"rapidEx/internal/domain/stock"
	"rapidEx/internal/storage"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type StockModify struct {
	Ticker string            `json:"ticker"`
	Price  float64           `json:"price"`
	Buy    map[string]string `json:"stockBookBuy"`
	Sell   map[string]string `json:"stockBookSell"`
}

func NewStockMapString(s *stock.Stock) *StockModify {
	var sMap StockModify
	s.Lock.Lock()
	sMap.Ticker = s.Ticker
	sMap.Price = s.Price
	sMap.Buy = mapFloatToString(s.Stockbook.Buy)
	sMap.Sell = mapFloatToString(s.Stockbook.Sell)
	s.Lock.Unlock()
	return &sMap
}

func (s StockModify) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func UnmarshalBinary(data []byte) (*StockModify, error) {
	var s StockModify
	err := json.Unmarshal(data, &s)
	return &s, err
}

type Repository interface {
	Set(ctx context.Context, stock *stock.Stock) error
	Stock(ctx context.Context, ticker string) (*stock.Stock, error)
	Stocks(ctx context.Context) ([]*stock.Stock, error)
	Del(ctx context.Context, ticker string) error
}

type rsClient struct {
	rc *redis.Client
}

func (r *rsClient) Set(ctx context.Context, stock *stock.Stock) error {
	stockMapString := NewStockMapString(stock)
	status := r.rc.Set(ctx, stock.Ticker, stockMapString, time.Second*1000)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (r *rsClient) Stock(ctx context.Context, ticker string) (*stock.Stock, error) {
	rStockMapString := r.rc.Get(ctx, ticker)
	switch {
	case errors.Is(rStockMapString.Err(), redis.Nil):
		return nil, storage.ErrStockNotFound
	case rStockMapString.Err() != nil:
		return nil, rStockMapString.Err()
	}

	var s stock.Stock

	sStockWrap, err := rStockMapString.Result()
	if err != nil {
		return nil, err
	}

	stockMapString, err := UnmarshalBinary([]byte(sStockWrap))
	if err != nil {
		return nil, err
	}

	mBuy, err := mapStringToFloat(stockMapString.Buy)
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

func (r *rsClient) Stocks(ctx context.Context) ([]*stock.Stock, error) {
	panic("implement me !")
}

func (r *rsClient) Del(ctx context.Context, ticker string) error {
	intCmd := r.rc.Del(ctx, ticker)

	if intCmd.Err() != nil {
		return intCmd.Err()
	}

	return nil
}

func NewStockRepository(rc *redis.Client) Repository {
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
