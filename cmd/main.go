package main

import (
	"log"

	"rapidEx/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	app.ListenAndServe()
}