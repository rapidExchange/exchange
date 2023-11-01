package stockPriceProcessor

import (
	"math"
	"errors"
	stockBook "rapidEx/internal/domain/stock-book"
)

type stockPriceProcessor struct {
}

func New() *stockPriceProcessor {
	return &stockPriceProcessor{}
}

func meanWeight(book map[float64]float64) float64 { // returns the meanweight price
	priWeight, suquant := 0.0, 0.0
	for pri, quant := range book {
		priWeight += pri * quant
		suquant += quant
	}
	meanw := priWeight / suquant
	return meanw
}

//TODO: refactor
func (proc *stockPriceProcessor) MeanWeight(stockBook stockBook.StockBook) (float64, error) {
	buyWeight := meanWeight(stockBook.Buy)
	sellWeight := meanWeight(stockBook.Sell)
	if math.IsNaN(buyWeight) {
		if math.IsNaN(sellWeight) {
			return 0.0, errors.New("NaN price")
		}
		return sellWeight, nil
	}
	if math.IsNaN(sellWeight) {
		if math.IsNaN(buyWeight) {
			return 0.0, errors.New("NaN price")
		}
		return buyWeight, nil
	}
	return (sellWeight + buyWeight)/2, nil
}

func (proc *stockPriceProcessor) Round(x float64, prec int) float64 { // round to a floor
	pow := math.Pow10(prec)
	rounded := math.Floor(x * pow)
	return rounded / pow
}

func (proc *stockPriceProcessor) PreciseAs(part float64) int { // returns the number of significant decimal places
	k := 0
	if part < 0.99 {
		for int(part)%10 == 0 {
			k++
			part *= 10
		}
	}
	return k
}

// func (proc *stockPriceProcessor) UpdPrice(stb *stockBook.StockBook, quantmin uint64) float64 { // returns a rounded actual price
// 	comb := make(map[float64]float64)
// 	for k, v := range stb.Buy {
// 		comb[k] = v
// 	}
// 	for k, v := range stb.Sell {
// 		comb[k] = v
// 	}
// 	newprice := proc.MeanWeight(comb)
// 	minval := 1 / float64(quantmin)
// 	return proc.Round(newprice, proc.PreciseAs(minval))
// }
