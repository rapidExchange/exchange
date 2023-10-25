package stockPriceProcessor

import "math"

type stockPriceProcessor struct {
}

func New() *stockPriceProcessor {
	return &stockPriceProcessor{}
}

func (proc *stockPriceProcessor) UpdPrice(book map[float64]float64) float64 { // returns the meanweight price
	priWeight, suquant := 0.0, 0.0
	for pri, quant := range book {
		priWeight += pri * quant
		suquant += quant
	}
	meanw := priWeight / suquant
	return meanw
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
