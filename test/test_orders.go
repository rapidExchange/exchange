package test

import (
	"context"
	"rapidEx/internal/domain/order"
	"rapidEx/internal/redis-connect"
	"testing"
)

//TODO: finish test func
func OrderTest(t *testing.T) {
	order1 := order.New("btc/usdt", "user_1", 1.41, 36000.3)

	rc, err := redisconnect.SetRedisConn()
	if err != nil {
		t.Error(err)
	}

	orderRepository := order.NewRepository(rc)

	err = orderRepository.Set(context.Background(), *order1)
	if err != nil {
		t.Error(err)
	}
} 