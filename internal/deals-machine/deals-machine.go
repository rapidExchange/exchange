package dealsprocessor

import (
	"context"
	"fmt"
	"log"
	"rapidEx/internal/domain/order"
	"rapidEx/internal/domain/stock"
)

type dealsMachine struct {
	ctx           context.Context
	ordersStorage OrdersStorage
	stockProvider stockProvider
}

type stockProvider interface {
	Stock(ctx context.Context, ticker string) (*stock.Stock, error)
}

type OrdersStorage interface {
	All(ctx context.Context, s *stock.Stock) ([]*order.Order, error)
	Del(ctx context.Context, order *order.Order)
}

// TODO: add order in order history & update user balance-sheet
// TODO send order to mysql order table
// TODO: rewrite function to worker with stock and orders in arguments
func (d *dealsMachine) Do(s *stock.Stock) {
	orders, err := d.ordersStorage.All(d.ctx, s)
	if err != nil {
		panic(err)
	}
	for _, order := range orders {

		if order.Status == "processing" {
			ok := processOrder(order, s)
			if ok {
				log.Printf("order %s from user %s was successfully done", order.OrderUUID.String(), order.Email)
				if err != nil {
					panic(fmt.Sprint("deals-processor: ", err))
				}
				order.Status = "completed"
				d.ordersStorage.Del(context.Background(), order)
			}
		}

	}
}

func processOrder(order *order.Order, stock *stock.Stock) bool {
	switch {
	case order.Type == "b" && stock.Price <= order.Price:
		return processBuyOrder(order, stock)
	case order.Type == "s" && stock.Price >= order.Price:
		return processSellOrder(order, stock)
	}
	return false
}

func processBuyOrder(order *order.Order, stock *stock.Stock) bool {
	for price, quantity := range stock.Stockbook.Sell {
		if order.Price >= price && order.Quantity <= quantity {
			stock.Stockbook.Sell[price] -= order.Quantity
			return true
		}
	}
	return false
}

func processSellOrder(order *order.Order, stock *stock.Stock) bool {
	for price, quantity := range stock.Stockbook.Buy {
		if order.Price <= price && order.Quantity <= quantity {
			stock.Stockbook.Buy[price] -= order.Quantity
			return true
		}
	}
	return false
}

func New(ctx context.Context, ordersStorage OrdersStorage, stockProvider stockProvider) *dealsMachine {
	return &dealsMachine{ctx: ctx, ordersStorage: ordersStorage, stockProvider: stockProvider}
}
