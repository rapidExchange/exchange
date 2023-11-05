package order

import (
	"context"
	"strconv"
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
	status := r.rc.Set(ctx, order.Ticker, order, 1)
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