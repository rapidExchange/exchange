package stockBook

type StockBook struct {
	buy *map[float64]float64
	sell *map[float64]float64
}

func New() *StockBook {
	buy := make(map[float64]float64)
	sell := make(map[float64]float64)
	return &StockBook{
		buy: &buy,
		sell: &sell,
	}
}