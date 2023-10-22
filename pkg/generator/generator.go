package generator

import (
	"math"
	"math/rand"
)

func generate(cPrice float64) (volume float64, price float64) {
	volume = float64(rand.Int31n(5000-10)+10) + rand.Float64()
	price = math.Abs(rand.NormFloat64()*0.05*cPrice + cPrice) // editable compression parameter
	return
}
