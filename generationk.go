package generationk

import (
	"fmt"
	inter "generationk/internal"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

//PutEvent starts a backtest with the information in context
func PutEvent(c *inter.Context, data chan inter.Event) {

}

func MakeOrder(ctx *inter.Context, ordertype inter.OrderType, asset *inter.Asset, time time.Time, amount float64) {
	log.WithFields(log.Fields{
		"Asset":  asset.Name,
		"Time":   time,
		"Amount": amount,
	}).Info("GENERATIONK>MAKE ORDER>")
	ctx.OrderChannel() <- inter.Order{
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

func RunLive(ctx *inter.Context) {

	//Initialize the strategy
	//ctx.Strategy[0].Setup(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go run(ctx, &wg)
	wg.Wait()
	//go dataManager.Start()
}

//Run starts a backtest with the information in context
func run(ctx *inter.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var o sync.Once
	//var init bool
	//var once sync.Once
	//Quit:
	for {
		log.Debug("GENERATIONK> MAIN LOOP")
		log.WithFields(log.Fields{
			"Number of items": len(ctx.OrderChannel()),
		}).Debug("GENERATIONK>ORDER CHANNEL")
		select {
		case orderEvent := <-ctx.OrderChannel():
			switch orderEvent.(type) {

			case inter.Order:
				log.Debug("GENERATIONK>ORDERCHANNEL> ORDER EVENT PICKED OFF QUEUE")
				ctx.Broker.PlaceOrder(orderEvent.(inter.Order))

			case inter.Fill:
				log.Debug("GENERATIONK>ORDERCHANNEL> FILL EVENT PICKED OFF QUEUE")
				ctx.Portfolio.Fill(orderEvent.(inter.Fill))
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
				case inter.Tick:
					log.Debug("GENERATIONK>EVENTCHANNEL> TICK EVENT PICKED OFF QUEUE")
					//fmt.Printerln("Processing tick data")
					log.Debug("GENERATIONK>EVENTCHANNEL> Leting strategy know")
					for i := range ctx.Strategy {
						ctx.Strategy[i].Tick(ctx)
					}
				case inter.DataEvent:
					if ctx.EndDate.After(event.(inter.DataEvent).Ohlc.Time) {
						log.Debug("GENERATIONK>EVENTCHANNEL> Ohlc.Time is after the back test end date")
						ctx.EventChannel() <- inter.Quit{}
						break
					}
					ctx.K++
					log.Debug("GENERATIONK>EVENTCHANNEL> DATAEVENT EVENT PICKED OFF QUEUE")
					//Add data to asset
					if _, ok := ctx.AssetMap[event.(inter.DataEvent).Name]; ok {
						log.WithFields(log.Fields{
							"(inter.DataEvent).Name": event.(inter.DataEvent).Name,
						}).Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> EXISTS IN MAP")
						//do something here
						//ctx.AssetMap[event.(inter.DataEvent).Name].Ohlc = prepend(ctx.AssetMap[event.(inter.DataEvent).Name].Ohlc, event.(inter.DataEvent).Ohlc)
						ctx.AssetMap[event.(inter.DataEvent).Name].Update(event.(inter.DataEvent).Ohlc)

						//resize(ctx.AssetMap[event.(inter.DataEvent).Name].Ohlc, ctx.Strategy[0].GetInitPeriod())
					} else {
						log.Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> CREATING ASSET AND ADDING TO MAP")
						//var asset inter.Asset
						//asset.Name = event.(inter.DataEvent).Name
						//asset.Update(event.(inter.DataEvent).Ohlc)
						asset := inter.NewAsset(event.(inter.DataEvent).Name, event.(inter.DataEvent).Ohlc)
						//asset.Ohlc = append(asset.Ohlc, )
						ctx.AssetMap[event.(inter.DataEvent).Name] = asset
					}
					//Run only once at some point

					o.Do(func() {
						ctx.Strategy[0].Setup(ctx)
						log.Debug("GENERATIONK>RUN ONCE")
					})
					//Run setup after initperiod is finished
					if ctx.K < ctx.Strategy[0].GetInitPeriod() {

						log.Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> Initializing strategy failed")
						break

					} else {

						log.Debug("GENERATIONK>EVENTCHANNEL> Updating indicators data")
						updateIndicators(ctx, event.(inter.DataEvent))

						log.Debug("GENERATIONK>EVENTCHANNEL> Leting strategy know")
						for i := range ctx.Strategy {
							//ctx.Strategy[i].Update(ctx)
							ctx.Strategy[i].Tick(ctx)
						}
					}
				case inter.Quit:
					log.Debug("GENERATIONK>EVENTCHANNEL> QUIT EVENT PICKED OFF QUEUE")
					close(ctx.OrderChannel())
					close(ctx.EventChannel())
					return
					//break Quit
				}
				//default:
				//ctx.IncOneDay()
				//ctx.EventChannel() <- inter.Tick{}
				//log.Debug("GENERATIONK>EVENTCHANNEL> NEW TICK ON QUEUE")
				//log.Debug("GENERATIONK>EVENTCHANNEL> DEFAULT")
				/*_, error := shiftData(ctx)
				if error != nil {
					data <- Quit{}
				}*/
			}
		}
	}
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func updateIndicators(ctx *inter.Context, dataEvent inter.DataEvent) {
	log.Debug("ctx.AssetIndicatorMap[dataEvent.Name]: ", len(ctx.AssetIndicatorMap[dataEvent.Name]))
	for k := range ctx.AssetIndicatorMap[dataEvent.Name] {
		indicator := (*ctx.AssetIndicatorMap[dataEvent.Name][k])
		data := ctx.AssetMap[dataEvent.Name].CloseArray()
		if len(data) > 0 {
			period := Min(len(ctx.AssetMap[dataEvent.Name].CloseArray()), indicator.GetPeriod())
			dataWindow := make([]float64, period)
			copy(dataWindow, data[:period])
			log.WithFields(log.Fields{
				"len(dataWindow)": len(dataWindow),
				"dataWindow":      dataWindow,
			}).Debug("GENERATIONK>UPDATE INDICATORS>")
			indicator.Update(dataWindow)
		}
		log.Debug("K: ", k)
	}
}

type EndOfDataError struct {
	Description string
}

func (e *EndOfDataError) Error() string {
	return fmt.Sprintf("End of data: %s", e.Description)
}

/*func shiftData(ctx *interernal.Context) (bool, error) {
	for i := range ctx.Asset {
		error := ctx.Asset[i].Shift(ctx.Time())
		if error != nil {
			fmt.Printerf("Data error <==============")
			return false, error
		}
	}
	return true, nil
}*/

func typeOfEvent(tst interface{}) {

}

//PutData starts a backtest with the information in context
func PutData(c *inter.Context, data chan inter.OHLC) {
	for _, asset := range c.Asset {
		for _, ohlc := range asset.Ohlc {
			data <- ohlc
		}
	}
}

//GetData starts a backtest with the information in context
func GetData(c *inter.Context, data chan inter.OHLC) {
	for _, asset := range c.Asset {
		for _, ohlc := range asset.Ohlc {
			data <- ohlc
		}
	}
}
