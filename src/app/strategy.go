package main

import (
	"fmt"
	"indicators"
)

//Strategy strategy
type Strategy struct {
}

func (m *Strategy) init(context *Context) {
	fmt.Printf("Init strategy\n")
	var ma indicators.Average
	floats := ma.Sma([]float64{1.0, 2.0, 3.0}, 7)
	fmt.Printf("MA: %f", floats)
}

func (m *Strategy) tick(ohlc OHLC) {
	fmt.Printf("Close: %f", ohlc.close)
}
