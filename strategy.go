package main

import (
	"fmt"
	ind "generationk/indicators"
)

//Strategy strategy
type Strategy struct {
}

func (m *Strategy) init(context *Context) {
	fmt.Printf("Init strategy\n")
	fmt.Printf("MA: %f", ind.Sma(7, context.asset.getCloseArray()))
}

func (m *Strategy) tick(ohlc OHLC) {
	fmt.Printf("Close: %f", ohlc.close)
}
