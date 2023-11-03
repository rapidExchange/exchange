package tickerstorage

import "sync"

var instance *tickerStorage = nil
var once sync.Once

type TickerStorage interface {
	TickerAppend(ticker string)
}

//tickerStorage is an implementation of singleton pattern for store all stock tickers
type tickerStorage struct{
	Tickers map[string]struct{}
	sync.Mutex
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