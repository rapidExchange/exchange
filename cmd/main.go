package main

import (
	"log"

	"rapidEx/internal/app"
	dealsprocessor "rapidEx/internal/deals-processor"
	"rapidEx/internal/generator"
	stockPriceProcessor "rapidEx/internal/stock-price-processor"
)

func main() {
	gen := generator.New()
	dealsprocessor := dealsprocessor.New()
	stockPriceProcessor := stockPriceProcessor.New()
	app, err := app.New(gen, dealsprocessor, stockPriceProcessor)
	if err != nil {
		log.Fatal(err)
	}
	go app.Do()
	app.ListenAndServe()
}
