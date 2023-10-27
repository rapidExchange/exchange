package stockBook

type StockBook struct {
	Buy  *map[float64]float64
	Sell *map[float64]float64
}

func New() *StockBook {
	buy := make(map[float64]float64)
	sell := make(map[float64]float64)
	return &StockBook{
		Buy:  &buy,
		Sell: &sell,
	}
}
