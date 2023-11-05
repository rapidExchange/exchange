package tickerstorage

import "sync"

var instance *tickerStorage = nil
var once sync.Once

type TickerStorage interface {
	TickerAppend(ticker string, precision int)
	GetTickers() []string
}

// tickerStorage is an implementation of singleton pattern for store all stock tickers
// the key of map is a ticker, value is precision(used for round prices and orders)
type tickerStorage struct {
	TickersPrecision map[string]int
	sync.RWMutex
}

func GetInstanse() TickerStorage {
	once.Do(func() {
		instance = &tickerStorage{TickersPrecision: make(map[string]int)}
	})
	return instance
}

func (t *tickerStorage) TickerAppend(ticker string, precision int) {
	t.Lock()
	defer t.Unlock()
	t.TickersPrecision[ticker] = precision
}

func (t *tickerStorage) GetTickers() []string {
	tickers := make([]string, 0)
	t.RLock()
	for k := range t.TickersPrecision {
		tickers = append(tickers, k)
	}
	t.RUnlock()
	return tickers
}
