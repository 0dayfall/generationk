package generationk

import (
	"fmt"
	int "generationk/internal"
	"time"

	log "github.com/sirupsen/logrus"
)

//PutEvent starts a backtest with the information in context
func PutEvent(c *int.Context, data chan int.Event) {

}

func MakeOrder(ctx *int.Context, ordertype int.OrderType, asset *int.Asset, time time.Time, amount float64) {
	log.Debug("GENERATIONK> Puting order on the orderchannel")
	ctx.OrderChannel() <- int.Order{
		Ordertype: ordertype,
		Asset:     asset,
		Time:      time,
		Amount:    amount,
	}
}

func Signal() {

}

func NewDataManager() {

}

func RunBacktest(ctx *int.Context, dataManager int.DataManager) {

	//Initialize the strategy
	ctx.Strategy[0].Setup(ctx)
	go run(ctx)
	go dataManager.Start()
}

//Run starts a backtest with the information in context
func run(ctx *int.Context) {
	for {
		log.Debug("GENERATIONK> MAIN LOOP")
		log.WithFields(log.Fields{
			"Number of items": len(ctx.OrderChannel()),
		}).Debug("GENERATIONK>ORDER CHANNEL")
		select {
		case orderEvent := <-ctx.OrderChannel():
			switch orderEvent.(type) {
			case int.Order:
				log.Debug("GENERATIONK>ORDERCHANNEL> ORDER EVENT PICKED OFF QUEUE")
				ctx.Broker.PlaceOrder(orderEvent.(int.Order))
			case int.Fill:
				log.Debug("GENERATIONK>ORDERCHANNEL> FILL EVENT PICKED OFF QUEUE")
				ctx.Portfolio.Fill(orderEvent.(int.Fill))
				log.Debug("GENERATIONK>ORDERCHANNEL> GIVING NOTICE TO STRATEGY")
				for i := range ctx.Strategy {
					ctx.Strategy[i].OrderEvent(ctx)
				}
			default:
				log.WithFields(log.Fields{
					"event": orderEvent,
				}).Debug("GENERATIONK>ORDERCHANNEL> UNKNOWN EVENT")
				fmt.Printf("%s", orderEvent)
			}
		default:
			log.Debug("GENERATIONK>ORDERCHANNEL> EMPTY")
			log.WithFields(log.Fields{
				"Number of items": len(ctx.EventChannel()),
			}).Debug("GENERATIONK>EVENTCHANNEL>")
			select {
			case event := <-ctx.EventChannel():
				switch event.(type) {
				case int.Tick:
					log.Debug("GENERATIONK>EVENTCHANNEL> TICK EVENT PICKED OFF QUEUE")
					//fmt.Println("Processing tick data")
					log.Debug("GENERATIONK>EVENTCHANNEL> Leting strategy know")
					for i := range ctx.Strategy {
						ctx.Strategy[i].Tick(ctx)
					}
				case int.Quit:
					log.Debug("GENERATIONK>EVENTCHANNEL> QUIT EVENT PICKED OFF QUEUE")
					close(ctx.EventChannel())
					break
				}
			default:
				ctx.IncOneDay()
				ctx.EventChannel() <- int.Tick{}
				log.Debug("GENERATIONK>EVENTCHANNEL> NEW TICK ON QUEUE")
				/*_, error := shiftData(ctx)
				if error != nil {
					data <- Quit{}
				}*/
			}
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
