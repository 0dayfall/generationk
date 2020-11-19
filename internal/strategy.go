package internal

import (
	"fmt"
	"generationk/indicators"
)

//Strategy strategy
type Strategy struct {
}

//Init is used to start the strategy
func (m *Strategy) Init(context *Context) {
	fmt.Printf("Init strategy\n")
	var ma indicators.Average
	floats := ma.Sma([]float64{1.0, 2.0, 3.0}, 7)
	fmt.Printf("MA: %f", floats)
}

//Tick gets called when new data is available
func (m *Strategy) Tick(ohlc OHLC) {
	fmt.Printf("Close: %f", ohlc.close)
}
