package generationk

import (
	"fmt"
	internal "generationk/internal"
)

//Event type
type Event interface {
	Handle()
}

//Tick event type
type Tick struct{}

func (t Tick) Handle() {}

//Signal event type
type Signal struct{}

func (s Signal) Handle() {}

//Order event type
type Order struct{}

func (o Order) Handle() {}

//Fill event type
type Fill struct{}

func (f Fill) Handle() {}

//Data event type
type Data struct{}

func (d Data) Handle() {}

//PutEvent starts a backtest with the information in context
func PutEvent(c *internal.Context, data chan Event) {

}

func RunBacktest(ctx *internal.Context) {
	eventChannel := make(chan Event, 1)
	run(ctx, eventChannel)
}

//Run starts a backtest with the information in context
func run(ctx *internal.Context, data chan Event) {
	for {
		select {
		case event := <-data:
			switch event.(type) {
			case Tick:
				fmt.Println("Processing tick data")
				for _, strat := range ctx.Strategy {
					strat.Tick(ctx)
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
			ctx.IncTime()
			data <- Tick{}
		}
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
