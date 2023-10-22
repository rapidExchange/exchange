package main

import (
	"fmt"
	"math"
	"math/rand"
)

func generate(cPrice float64) (volume float64, price float64) {
	volume = float64(rand.Int31n(5000 -10) + 10) + rand.Float64()
	min := int32(0.7*cPrice)
	max := int32(1.3*cPrice)
	price = math.Abs(rand.NormFloat64()) + float64(rand.Int31n(max - min) + min) - 0.3
	return
}

func main() {
	sum := 0.0
	for i := 0; i < 1000000; i++ {
		_, pri := generate(10000)
		sum += pri
	}
	fmt.Println(sum/1000000)
}