package generationk

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func MakeOrder(ctx *Context, ordertype OrderType, asset *Asset, time time.Time, amount float64) {
	log.Debug("GENERATIONK>makeOrder()")
	makeOrder(ctx, ordertype, asset, time, amount)
}

func makeOrder(ctx *Context, ordertype OrderType, asset *Asset, time time.Time, amount float64) {
	log.WithFields(log.Fields{
		"Asset":  asset.Name,
		"Time":   time,
		"Amount": amount,
	}).Info("GENERATIONK>MAKE ORDER>")

	ctx.OrderChannel() <- Order{
		Ordertype: ordertype,
		Asset:     asset,
		Time:      time,
		Amount:    amount,
	}
}

func RunEventBased(ctx *Context) {
	//Initialize the strategy
	var wg sync.WaitGroup
	wg.Add(1)
	go run(ctx, &wg)
	wg.Wait()
}

//Run starts a backtest with the information in context
func run(ctx *Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var o sync.Once
	for {
		select {
		case orderEvent := <-ctx.OrderChannel():
			log.WithFields(log.Fields{
				"Number of items": len(ctx.OrderChannel()),
			}).Debug("GENERATIONK>ORDER EVENT PICKED OFF QUEUE")

			switch orderEvent.(type) {

			case Order:
				log.Debug("GENERATIONK>ORDERCHANNEL>ORDER> EVENT PICKED OFF QUEUE")
				//				go func() {
				log.Debug("GENERATIONK>ORDERCHANNEL>ORDER>PUTING SUBMITTED ON QUEUE")
				ctx.OrderChannel() <- Submitted{}
				//				}()
				ctx.Broker.PlaceOrder(orderEvent.(Order))
			case Submitted:
				log.Debug("GENERATIONK>ORDERCHANNEL>SUBMIT> EVENT PICKED OFF QUEUE")

				//for i := range ctx.Strategy {
				ctx.Strategy[0].OrderEvent(orderEvent)
				//}

			case Accepted:
				log.Debug("GENERATIONK>ORDERCHANNEL>ACCEPT> EVENT PICKED OFF QUEUE")

				//for i := range ctx.Strategy {
				ctx.Strategy[0].OrderEvent(orderEvent)
				//}

			case PartialFill:
				log.Debug("GENERATIONK>ORDERCHANNEL>PARTIALFILL> EVENT PICKED OFF QUEUE")

				//for i := range ctx.Strategy {
				ctx.Strategy[0].OrderEvent(orderEvent)
				//}

			case Fill:
				log.Debug("GENERATIONK>ORDERCHANNEL>FILL> EVENT PICKED OFF QUEUE")

				//for i := range ctx.Strategy {
				ctx.Strategy[0].OrderEvent(orderEvent)
				//}

			case Rejected:
				log.Debug("GENERATIONK>ORDERCHANNEL>REJECTED> EVENT PICKED OFF QUEUE")

				//for i := range ctx.Strategy {
				ctx.Strategy[0].OrderEvent(orderEvent)
				//}
			default:
				log.WithFields(log.Fields{
					"event": orderEvent,
				}).Debug("GENERATIONK>ORDERCHANNEL> DEFAULTS - ORDERCHANNEL EMPTY")
			}
		default:

			select {
			case event := <-ctx.EventChannel():

				log.WithFields(log.Fields{
					"Number of items": len(ctx.EventChannel()),
				}).Debug("GENERATIONK>DATA EVENT PICKED OFF QUEUE")

				switch event.(type) {
				case Tick:

					log.Debug("GENERATIONK>EVENTCHANNEL> TICK EVENT PICKED OFF QUEUE")
					//fmt.Printerln("Processing tick data")
					log.Debug("GENERATIONK>EVENTCHANNEL> Leting strategy know")
					for i := range ctx.Strategy {
						ctx.Strategy[i].Tick(ctx)
					}

				case DataEvent:

					if ctx.EndDate.After(event.(DataEvent).Ohlc.Time) {
						log.Debug("GENERATIONK>EVENTCHANNEL> Ohlc.Time is after the back test end date")
						ctx.EventChannel() <- Quit{}
						break
					}
					ctx.K++

					//Add data to asset
					if _, ok := ctx.AssetMap[event.(DataEvent).Name]; !ok {
						log.Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> CREATING ASSET AND ADDING TO MAP")
						asset := NewAsset(event.(DataEvent).Name, event.(DataEvent).Ohlc)
						ctx.AssetMap[event.(DataEvent).Name] = asset
					}

					log.WithFields(log.Fields{
						"(DataEvent).Name": event.(DataEvent).Name,
					}).Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> EXISTS IN MAP")

					//ctx.AssetMap[event.(DataEvent).Name].Ohlc = prepend(ctx.AssetMap[event.(DataEvent).Name].Ohlc, event.(DataEvent).Ohlc)
					ctx.AssetMap[event.(DataEvent).Name].Update(event.(DataEvent).Ohlc)

					//Run only once to setup indicators
					o.Do(func() {
						log.Debug("GENERATIONK>RUN ONCE")
						ctx.Strategy[0].Setup(ctx)
						log.WithFields(log.Fields{
							"strategy": ctx.Strategy[0],
						}).Debug("Strategy")
					})

					//Run setup after initperiod is finished
					if ctx.K < ctx.GetInitPeriod() {

						log.Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> Initializing strategy failed")
						break

					} else {

						log.Debug("GENERATIONK>EVENTCHANNEL> Updating indicators data")
						updateIndicators(ctx, event.(DataEvent))

						log.Debug("GENERATIONK>EVENTCHANNEL> Leting strategy know")
						ctx.Strategy[0].Tick(ctx)
					}
				case Quit:
					log.Debug("GENERATIONK>EVENTCHANNEL> QUIT EVENT PICKED OFF QUEUE")
					close(ctx.OrderChannel())
					close(ctx.EventChannel())
					return

				}
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

func updateIndicators(ctx *Context, dataEvent DataEvent) {
	log.Debug("ctx.AssetIndicatorMap[dataEvent.Name]: ", len(ctx.AssetIndicatorMap[dataEvent.Name]))

	//If the asset has no data so far ther is no point in doing this
	data := ctx.AssetMap[dataEvent.Name].CloseArray()
	if len(data) < 1 {
		return
	}

	for k := range ctx.AssetIndicatorMap[dataEvent.Name] {

		indicator := (*ctx.AssetIndicatorMap[dataEvent.Name][k])

		//Copy either the data we have available or period much to the indicator
		period := Min(len(ctx.AssetMap[dataEvent.Name].CloseArray()), indicator.GetPeriod())
		dataWindow := make([]float64, period)
		copy(dataWindow, data[:period])

		log.WithFields(log.Fields{
			"len(dataWindow)": len(dataWindow),
			"dataWindow":      dataWindow,
		}).Debug("GENERATIONK>UPDATE INDICATORS>")

		//Update the indicator with new data
		indicator.Update(dataWindow)

		log.Debug("K: ", k)
	}
}

type EndOfDataError struct {
	Description string
}

func (e *EndOfDataError) Error() string {
	return fmt.Sprintf("End of data: %s", e.Description)
}
