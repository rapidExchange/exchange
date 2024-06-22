package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"rapidEx/internal/app"
	dealsMachine "rapidEx/internal/deals-machine"
	redisconnect "rapidEx/internal/redis"
	orderrepository "rapidEx/internal/repositories/order-repository"
	stockrepository "rapidEx/internal/repositories/stock-repository"
	"rapidEx/internal/services/order"
	"rapidEx/internal/services/stock"
	stockPriceProcessor "rapidEx/internal/stock-price-processor"
)

func main() {
	rc := redisconnect.MustConnect()
	orderRepository := orderrepository.NewOrderRepository(rc)
	orderMonitor := order.New(slog.New(slog.NewTextHandler(os.Stdout, nil)), orderRepository, orderRepository, orderRepository)
	stockrepository := stockrepository.NewStockRepository(rc)
	stockMonitor := stock.New(slog.New(slog.NewTextHandler(os.Stdout, nil)), stockrepository, stockrepository, nil)
	dealsMachine := dealsMachine.New(context.Background(), orderMonitor, stockMonitor)
	stockPriceProcessor := stockPriceProcessor.New()
	app, err := app.New(dealsMachine, stockPriceProcessor)
	if err != nil {
		log.Fatal(err)
	}
	go app.Do()
	app.ListenAndServe()
}
