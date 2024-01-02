package dealsprocessor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"rapidEx/internal/domain/order"
	redisconnect "rapidEx/internal/redis"
	orderrepository "rapidEx/internal/repositories/order-repository"
	stockrepository "rapidEx/internal/repositories/stock-repository"
	"rapidEx/internal/services/stock"
	stockDomain "rapidEx/internal/domain/stock"
	"rapidEx/internal/storage"
)

type dealsProcessor struct {
}

// TODO: add order in order history & update user balance-sheet
func (d *dealsProcessor) Do() {
	orders, err := getAllOrders()
	stockRepository := stockrepository.NewStockRepository(redisconnect.MustConnect())
	if err != nil {
		panic(err)
	}
	stockMonitor := stock.New(&slog.Logger{}, 
		stockRepository,
		stockRepository, nil)
	for _, order := range orders {

		stock, err := stockMonitor.Stock(context.Background(), order.Ticker)
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Printf("Stock: %s not found\n", order.Ticker)
			order.Status = "unable"
			continue
		}

		if order.Status == "processing" {
			ok := processOrder(order, stock)
			if ok {
				log.Printf("order %s from user %s was successfully done", order.OrderUUID.String(), order.Email)
				err = stockRepository.Set(context.Background(), stock)
				if err != nil {
					panic(fmt.Sprint("deals-processor: ", err))
				}
				orderRepository := orderrepository.NewOrderRepository(redisconnect.MustConnect())
				order.Status = "completed"
				//TODO send order to mysql order table
				orderRepository.Del(context.Background(), order)
			}
		}

	}
}

func processOrder(order *order.Order, stock *stockDomain.Stock) bool {
	switch {
	case order.Type == "b" && stock.Price <= order.Price:
		return processBuyOrder(order, stock)
	case order.Type == "s" && stock.Price >= order.Price:
		return processSellOrder(order, stock)
	}
	return false
}

func processBuyOrder(order *order.Order, stock *stockDomain.Stock) bool {
	for price, quantity := range stock.Stockbook.Sell {
		if order.Price >= price && order.Quantity <= quantity {
			stock.Stockbook.Sell[price] -= order.Quantity
			return true
		}
	}
	return false
}

func processSellOrder(order *order.Order, stock *stockDomain.Stock) bool {
	for price, quantity := range stock.Stockbook.Buy {
		if order.Price <= price && order.Quantity <= quantity {
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

func New() *dealsProcessor {
	return &dealsProcessor{}
}
