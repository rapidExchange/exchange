package stockBook

type StockBook struct {
	buy *map[string]int64
	sell *map[string]int64
}

func New() *StockBook {
	buy := make(map[string]int64)
	sell := make(map[string]int64)
	return &StockBook{
		buy: &buy,
		sell: &sell,
	}
}