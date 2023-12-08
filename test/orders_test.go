package test

import (
	"context"
	"rapidEx/internal/domain/order"
	redisconnect "rapidEx/internal/redis-connect"
	orderrepository "rapidEx/internal/repositories/order-repository"
	"testing"
)

func TestOrder(t *testing.T) {
	order1, err := order.New("btc/usdt", "user_1", "b", 1.41, 36000.3)
	if err != nil {
		t.Error(err)
	}

	rc, err := redisconnect.SetRedisConn()
	if err != nil {
		t.Error(err)
	}

	orderRepository := orderrepository.NewOrderRepository(rc)

	err = orderRepository.Set(context.Background(), order1)
	if err != nil {
		t.Error(err)
	}
	orders, err := orderRepository.GetAll(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range orders {
		if *order1 == *v {
			err = orderRepository.Del(context.Background(), order1)
			if err != nil {
				t.Error(err)
			}
		}
	}
}
