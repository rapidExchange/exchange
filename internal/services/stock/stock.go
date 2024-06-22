package stock

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"rapidEx/internal/domain/stock"
	"rapidEx/internal/storage"
)

var (
	ErrStockNotFound = errors.New("stock not found")
)

type StockProvider interface {
	Stock(ctx context.Context, ticker string) (*stock.Stock, error)
}

type StockSaver interface {
	Set(ctx context.Context, stock *stock.Stock) error
}

type PrecisionProvider interface {
	Precision(price float64) int64
}

type StockService struct {
	log               *slog.Logger
	stockProvider     StockProvider
	stockSaver        StockSaver
	precisionProvider PrecisionProvider
}

func New(log *slog.Logger, stockProvider StockProvider,
	stockSaver StockSaver,
	precisionProvider PrecisionProvider) *StockService {
	return &StockService{log: log,
		stockProvider:     stockProvider,
		stockSaver:        stockSaver,
		precisionProvider: precisionProvider}
}

func (s *StockService) Stock(ctx context.Context, ticker string) (*stock.Stock, error) {
	const op = "stockMonitor.Stock"
	log := s.log.With(slog.String("op", op),
		slog.String("ticker", ticker))
	stock, err := s.stockProvider.Stock(ctx, ticker)
	if err != nil {
		if errors.Is(err, storage.ErrStockNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrStockNotFound)
		}
		log.Warn(err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return stock, nil
}

func (s *StockService) Set(ctx context.Context, stock *stock.Stock) error {
	const op = "stockMonitor.Stock"
	log := s.log.With(slog.String("op", op),
		slog.String("ticker", stock.Ticker))
	err := s.stockSaver.Set(ctx, stock)
	if err != nil {
		log.Warn("failed to save stock")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
