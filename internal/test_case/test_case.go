package case1

import (
	"fmt"
	"rapidEx/internal/generator"
	stockPriceProcessor "rapidEx/internal/stock-price-processor"
	"rapidEx/internal/usecases/stock_usecases"
	"time"
)

func Case() {
	for {
		time.Sleep(time.Second * 2)
		s, err := stock_usecases.GetStock("btc/usdt")
		if err != nil {
			fmt.Println(err)
		} else {
			gen := generator.New()
			for i := 0; i < 10; i++ {
				gen.OrderGenerate(s)
			}
			sProcessor := stockPriceProcessor.New()
			price, err := sProcessor.MeanWeight(s.Stockbook)
			if err != nil {
				fmt.Println(err)
				return
			}
			s.Price = price
			fmt.Printf("New price btc/usdt: %.4f\n", s.Price)
			err = stock_usecases.SetStock(s.Ticker, s.Price)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
