package generationk

import (
	"fmt"
	internal "generationk/internal"
)

//Event type
type Event interface{
	Handle()
}

//Tick event type
type Tick struct {}
func(t Tick) Handle() {}

//Signal event type
type Signal struct {}
func(s Signal) Handle() {}

//Order event type
type Order struct {}
func(o Order) Handle() {}

//Fill event type
type Fill struct {}
func(f Fill) Handle() {}

//Data event type
type Data struct {}
func(d Data) Handle() {}

//PutEvent starts a backtest with the information in context
func PutEvent(c *internal.Context, data chan Event) {

}

//Run starts a backtest with the information in context
func Run(c *internal.Context, data chan Event) {
	select {
	case event := <-data:
		switch event.(type) {
		case Tick:
			fmt.Println("Processing tick data")
			for _, strat := range c.Strategy {
				strat.Tick(c)
			}
		case Signal:
			// here v has type S
		case Order:
			// here v has type S
		case Fill:
			// here v has type S
		case Data:
			// here v has type S
		default:
			// no match; here v has the same type as i
		}		
	default:
		fmt.Println("Increase time")
		c.IncTime()
		var tick Tick 
		data <- tick
		fmt.Printf("Len channel %d", len(data))
	}
}

func typeOfEvent(tst interface{}) {
	
}

//PutData starts a backtest with the information in context
func PutData(c *internal.Context, data chan internal.OHLC) {
	for _, asset := range c.Asset {
		for _, ohlc := range asset.Ohlc {
			data <- ohlc
		}
	}
}

//GetData starts a backtest with the information in context
func GetData(c *internal.Context, data chan internal.OHLC) {
	for _, asset := range c.Asset {
		for _, ohlc := range asset.Ohlc {
			data <- ohlc
		}
	}
}
