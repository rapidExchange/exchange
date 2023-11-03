package tickerstorage

import "sync"

var instance *tickerStorage = nil
var once sync.Once

type TickerStorage interface {
	TickerAppend(ticker string)
	GetTickers() []string
}

//tickerStorage is an implementation of singleton pattern for store all stock tickers
type tickerStorage struct{
	Tickers map[string]struct{}
	sync.RWMutex
}

func GetInstanse() TickerStorage {
	once.Do(func() {
		instance = &tickerStorage{Tickers: make(map[string]struct{})}
	})
	return instance
}

func (t *tickerStorage) TickerAppend(ticker string) {
	t.Lock()
	defer t.Unlock()
	t.Tickers[ticker] = struct{}{}
}

func (t *tickerStorage) GetTickers() []string {
	tickers := make([]string, 0)
	t.RLock()
	for k := range t.Tickers {
		tickers = append(tickers, k)
	}
	t.RUnlock()
	return tickers
}