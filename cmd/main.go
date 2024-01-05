package main

import (
	"log"

	"rapidEx/internal/app"
	dealsprocessor "rapidEx/internal/deals-processor"
	"rapidEx/internal/generator"
	stockPriceProcessor "rapidEx/internal/stock-price-processor"
	tickerstorage "rapidEx/internal/tickerStorage"
)

func main() {
	gen := generator.New()
	dealsprocessor := dealsprocessor.New()
	stockPriceProcessor := stockPriceProcessor.New()
	tickerstorage := tickerstorage.GetInstanse()
	app, err := app.New(gen, dealsprocessor, stockPriceProcessor, tickerstorage)
	if err != nil {
		log.Fatal(err)
	}
	app.Do()
	app.ListenAndServe()
}
