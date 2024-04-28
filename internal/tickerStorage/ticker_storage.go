package tickerstorage

import "sync"

var instance *tickerStorage = nil
var once sync.Once

type TickerStorage interface {
	TickerAppend(ticker string, precision int)
	GetTickers() []string
	Find(ticker string) bool
	Get(ticker string) (int, bool)
}

// tickerStorage is an implementation of singleton pattern for store all stock tickers
// the key of map is a ticker, value is precision(used for round prices and orders)
type tickerStorage struct {
	TickersPrecision map[string]int
	sync.Mutex
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
	t.Lock()
	defer t.Unlock()
	for k := range t.TickersPrecision {
		tickers = append(tickers, k)
	}
	return tickers
}

// Find returns true if provided ticker is in the stirage, otherwise returns false
func (t *tickerStorage) Find(ticker string) bool {
	tickers := t.GetTickers()
	for _, t:= range tickers {
		if t == ticker {
			return true
		}
	}
	return false
}

func(t *tickerStorage) Get(ticker string) (int, bool) {
	prec, ok := t.TickersPrecision[ticker]
	return prec, ok
}