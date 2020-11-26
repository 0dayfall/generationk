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

//Fill event type
type Quit struct{}

func (q Quit) Handle() {}

//PutEvent starts a backtest with the information in context
func PutEvent(c *internal.Context, data chan Event) {

}

func RunBacktest(ctx *internal.Context) {
	eventChannel := make(chan Event, 1)
	ctx.Strategy[0].Init(ctx)
	run(ctx, eventChannel)
}

//Run starts a backtest with the information in context
func run(ctx *internal.Context, data chan Event) {
	for {
		select {
		case event := <-data:
			switch event.(type) {
			case Tick:
				//fmt.Println("Processing tick data")
				for i := range ctx.Strategy {
					ctx.Strategy[i].Tick(ctx)
				}
			case Signal:
				// here v has type S
			case Order:
				// here v has type S
			case Fill:
				// here v has type S
			case Data:
				// here v has type S
			case Quit:
				close(data)
				break
			default:
				// no match; here v has the same type as i
			}
		default:
			ctx.IncOneDay()
			/*_, error := shiftData(ctx)
			if error != nil {
				data <- Quit{}
			}*/
			data <- Tick{}
		}
	}
}

type EndOfDataError struct {
	Description string
}

func (e *EndOfDataError) Error() string {
	return fmt.Sprintf("End of data: %s", e.Description)
}

/*func shiftData(ctx *internal.Context) (bool, error) {
	for i := range ctx.Asset {
		error := ctx.Asset[i].Shift(ctx.Time())
		if error != nil {
			fmt.Printf("Data error <==============")
			return false, error
		}
	}
	return true, nil
}*/

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
