package utils

import (
	"math"
)

func Round(x float64, prec int) float64 {
	pow := math.Pow10(prec)
	rounded := math.Floor(x * pow)
	return rounded / pow
}