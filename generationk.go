package generationk

import (
	"fmt"
	int "generationk/internal"

	log "github.com/sirupsen/logrus"
)

//PutEvent starts a backtest with the information in context
func PutEvent(c *int.Context, data chan int.Event) {

}

func RunBacktest(ctx *int.Context) {
	ctx.Strategy[0].Indicators(ctx)
	go run(ctx)
}

//Run starts a backtest with the information in context
func run(ctx *int.Context) {
	for {
		select {
		case event := <-ctx.EventChannel():
			switch event.(type) {
			case int.Signal:
				// here v has type S
			case int.Order:
				//go ctx.Broker.Order(int.Buy, )
			case int.Fill:
				log.Debug("Received FILL")
				for i := range ctx.Strategy {
					ctx.Strategy[i].OrderEvent(ctx)
				}
				// here v has type S
			case int.Data:
				// here v has type S
			case int.Quit:
				log.Debug("Received QUIT")
				close(ctx.EventChannel())
				break
			case int.Tick:
				log.Debug("Received TICK")
				//fmt.Println("Processing tick data")
				for i := range ctx.Strategy {
					ctx.Strategy[i].Orders(ctx)
				}
			default:
				// no match; here v has the same type as i
			}
		default:
			go ctx.IncOneDay()
			/*_, error := shiftData(ctx)
			if error != nil {
				data <- Quit{}
			}*/
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
func PutData(c *int.Context, data chan int.OHLC) {
	for _, asset := range c.Asset {
		for _, ohlc := range asset.Ohlc {
			data <- ohlc
		}
	}
}

//GetData starts a backtest with the information in context
func GetData(c *int.Context, data chan int.OHLC) {
	for _, asset := range c.Asset {
		for _, ohlc := range asset.Ohlc {
			data <- ohlc
		}
	}
}
