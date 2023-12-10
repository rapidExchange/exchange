package dealsprocessor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rapidEx/internal/domain/order"
	"rapidEx/internal/domain/stock"
	redisconnect "rapidEx/internal/redis-connect"
	orderrepository "rapidEx/internal/repositories/order-repository"
	stockrepository "rapidEx/internal/repositories/stock-repository"
)

type DealsProcessor struct {
}

// TODO: add order in order history & update user balance-sheet
func Do() {
	orders, err := getAllOrders()
	stockRepository := stockrepository.NewStockRepository(redisconnect.MustConnect())
	if err != nil {
		panic(err)
	}
	for _, order := range orders {

		stock, err := stockRepository.Get(context.Background(), order.Ticker)
		if errors.Is(err, stockrepository.ErrUserNotFound) {
			log.Printf("Stock: %s not found\n", order.Ticker)
		}

		if order.Status == "processing" {
			ok := processOrder(order, stock)
			if ok {
				log.Printf("order %s from user %s was successfully done", order.OrderUUID.String(), order.Email)
				err = stockRepository.Set(context.Background(), *stock)
				if err != nil {
					panic(fmt.Sprint("deals-processor: ", err))
				}
			}
		}

	}
}

func processOrder(order *order.Order, stock *stock.Stock) bool {
	switch {
	case order.Type == "buy" && stock.Price >= stock.Price:
		return processBuyOrder(order, stock)
	case order.Type == "sell" && stock.Price <= stock.Price:
		return processSellOrder(order, stock)
	}
	return false
}

func processBuyOrder(order *order.Order, stock *stock.Stock) bool {
	for price, quantity := range stock.Stockbook.Sell {
		if order.Price <= price && order.Quantity <= quantity {
			stock.Stockbook.Buy[price] -= order.Quantity
			return true
		}
	}
	return false
}

func processSellOrder(order *order.Order, stock *stock.Stock) bool {
	for price, quantity := range stock.Stockbook.Sell {
		if order.Price >= price && order.Quantity <= quantity {
			stock.Stockbook.Buy[price] -= order.Quantity
			return true
		}
	}
	return false
}

func getAllOrders() ([]*order.Order, error) {
	redisClient := redisconnect.MustConnect()
	orderRepository := orderrepository.NewOrderRepository(redisClient)
	return orderRepository.GetAll(context.Background())
}

func New() *DealsProcessor {
	return &DealsProcessor{}
}
