package generationk

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var o sync.Once

type GenkCallback interface {
	IsOwning(asset string) bool
	OrderSend(assetName string, ordertype OrderType, amount float64, qty int)
}

type generationK struct {
	market *Context
}

func NewGenerationK() *generationK {

	generationK := &generationK{
		market: newContext(),
	}
	return generationK
}

func (k *generationK) DataEvent(dataEvent Event) {

	log.WithFields(log.Fields{
		"Number of items": len(k.market.EventChannel()),
	}).Debug("GENERATIONK>DATA EVENT PICKED OFF QUEUE")

	k.market.datePointer = dataEvent.(DataEvent).Ohlc.Time

	if dataEvent.(DataEvent).Ohlc.Time.After(k.market.EndDate) {
		log.Debug("GENERATIONK>EVENTCHANNEL> Ohlc.Time is after the back test end date")
		//k.market.EventChannel() <- Quit{}
		return
	}
	k.market.K++

	//Add data to asset
	if _, ok := k.market.AssetMap[dataEvent.(DataEvent).Name]; !ok {
		log.Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> CREATING ASSET AND ADDING TO MAP")
		asset := NewAsset(dataEvent.(DataEvent).Name, dataEvent.(DataEvent).Ohlc)
		k.market.AssetMap[dataEvent.(DataEvent).Name] = asset
	}

	log.WithFields(log.Fields{
		"(DataEvent).Name": dataEvent.(DataEvent).Name,
	}).Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> EXISTS IN MAP")

	//k.market.AssetMap[event.(DataEvent).Name].Ohlc = prepend(k.market.AssetMap[event.(DataEvent).Name].Ohlc, event.(DataEvent).Ohlc)
	k.market.AssetMap[dataEvent.(DataEvent).Name].Update(dataEvent.(DataEvent).Ohlc)

	//Run only once to setup indicators
	o.Do(func() {
		log.Info("GENERATIONK>RUN ONCE")
		k.market.Strategy[0].Setup(k.market)
		log.WithFields(log.Fields{
			"strategy": k.market.Strategy[0],
		}).Debug("Strategy")
	})

	//Run setup after initperiod is finished
	if k.market.K < k.market.GetInitPeriod() {

		log.Info("GENERATIONK>EVENTCHANNEL>DATAEVENT> Initializing strategy failed")
		return

	} else {

		log.Info("GENERATIONK>EVENTCHANNEL> Updating indicators data")
		updateIndicators(k.market, dataEvent.(DataEvent))

		log.Info("GENERATIONK>EVENTCHANNEL> Leting strategy know")
		k.market.Strategy[0].Tick(k)
	}

	log.WithFields(log.Fields{
		"K: ": k.market.K,
	}).Info("K")

}

func (k *generationK) Init() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	lvl, ok := os.LookupEnv("LOG_LEVEL")

	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.FatalLevel
	}
	// set global log level
	logrus.SetLevel(ll)
}

func (k *generationK) AddAsset(asset *Asset) {
	k.market.AddAsset(asset)
}

func (k *generationK) AddStrategy(strat *Strategy) {
	k.market.AddStrategy(strat)
}

func (k *generationK) SetBalance(balance float64) {
	k.market.Portfolio.SetCash(balance)
}

func (k *generationK) AddStartDate(startDate time.Time) {
	k.market.AddStartDate(startDate)
}

func (k *generationK) AddEndDate(endDate time.Time) {
	k.market.AddEndDate(endDate)
}

func (k generationK) OrderSend(assetName string, ordertype OrderType, amount float64, qty int) {
	log.Debug("GENERATIONK>makeOrder()")
	orderSend(k.market, ordertype, assetName, k.market.datePointer, amount, qty)
}

func orderSend(ctx *Context, ordertype OrderType, assetName string, time time.Time, amount float64, qty int) {
	log.WithFields(log.Fields{
		"Asset":  assetName,
		"Time":   time,
		"Amount": amount,
		"Qty":    qty,
	}).Debug("GENERATIONK>MAKE ORDER>")

	orderStatus, _ := interface{}(ctx.Strategy[0]).(OrderStatus)

	ctx.Broker.SendOrder(
		Order{
			Ordertype: ordertype,
			Asset:     ctx.AssetMap[assetName],
			Time:      time,
			Amount:    amount,
			Qty:       qty,
		},
		orderStatus,
	)

}

//OwnPosition is used to find out if we have a holding in an asset
func (k generationK) IsOwning(name string) bool {
	return k.market.Portfolio.IsOwning(name)
}

func (k *generationK) RunEventBased() {
	//Initialize the strategy
	var wg sync.WaitGroup
	wg.Add(1)
	go run(k, &wg)
	wg.Wait()
}

//Run starts a backtest with the information in context
func run(k *generationK, wg *sync.WaitGroup) {
	defer wg.Done()
	var o sync.Once
	for {
		select {
		case orderEvent := <-k.market.OrderChannel():
			log.WithFields(log.Fields{
				"Number of items": len(k.market.OrderChannel()),
			}).Debug("GENERATIONK>ORDER EVENT PICKED OFF QUEUE")

			switch orderEvent.(type) {

			case Order:
				log.Debug("GENERATIONK>ORDERCHANNEL>ORDER> EVENT PICKED OFF QUEUE")
				//				go func() {
				log.Debug("GENERATIONK>ORDERCHANNEL>ORDER>PUTING SUBMITTED ON QUEUE")
				//k.market.OrderChannel() <- Submitted{}
				//				}()

			case Submitted:
				log.Debug("GENERATIONK>ORDERCHANNEL>SUBMIT> EVENT PICKED OFF QUEUE")

				//for i := range k.market.Strategy {
				//k.market.Strategy[0].OrderEvent(orderEvent)
				//}

			case Accepted:
				log.Debug("GENERATIONK>ORDERCHANNEL>ACCEPT> EVENT PICKED OFF QUEUE")

				//for i := range k.market.Strategy {
				//k.market.Strategy[0].OrderEvent(orderEvent)
				//}

			case PartialFill:
				log.Debug("GENERATIONK>ORDERCHANNEL>PARTIALFILL> EVENT PICKED OFF QUEUE")

				//for i := range k.market.Strategy {
				//k.market.Strategy[0].OrderEvent(orderEvent)
				//}

			case Fill:
				log.Debug("GENERATIONK>ORDERCHANNEL>FILL> EVENT PICKED OFF QUEUE")

				//for i := range k.market.Strategy {
				//k.market.Strategy[0].OrderEvent(orderEvent)
				//}

			case Rejected:
				log.Debug("GENERATIONK>ORDERCHANNEL>REJECTED> EVENT PICKED OFF QUEUE")

				//for i := range k.market.Strategy {
				//k.market.Strategy[0].OrderEvent(orderEvent)
				//}
			default:
				log.WithFields(log.Fields{
					"event": orderEvent,
				}).Debug("GENERATIONK>ORDERCHANNEL> DEFAULTS - ORDERCHANNEL EMPTY")
			}
		default:

			select {
			case event := <-k.market.EventChannel():

				log.WithFields(log.Fields{
					"Number of items": len(k.market.EventChannel()),
				}).Debug("GENERATIONK>DATA EVENT PICKED OFF QUEUE")

				switch event.(type) {
				case Tick:

					log.Debug("GENERATIONK>EVENTCHANNEL> TICK EVENT PICKED OFF QUEUE")
					//fmt.Printerln("Processing tick data")
					log.Debug("GENERATIONK>EVENTCHANNEL> Leting strategy know")
					for i := range k.market.Strategy {
						k.market.Strategy[i].Tick(k)
					}

				case DataEvent:
					k.market.datePointer = event.(DataEvent).Ohlc.Time

					if event.(DataEvent).Ohlc.Time.After(k.market.EndDate) {
						log.Debug("GENERATIONK>EVENTCHANNEL> Ohlc.Time is after the back test end date")
						k.market.EventChannel() <- Quit{}
						break
					}
					k.market.K++

					//Add data to asset
					if _, ok := k.market.AssetMap[event.(DataEvent).Name]; !ok {
						log.Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> CREATING ASSET AND ADDING TO MAP")
						asset := NewAsset(event.(DataEvent).Name, event.(DataEvent).Ohlc)
						k.market.AssetMap[event.(DataEvent).Name] = asset
					}

					log.WithFields(log.Fields{
						"(DataEvent).Name": event.(DataEvent).Name,
					}).Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> EXISTS IN MAP")

					//k.market.AssetMap[event.(DataEvent).Name].Ohlc = prepend(k.market.AssetMap[event.(DataEvent).Name].Ohlc, event.(DataEvent).Ohlc)
					k.market.AssetMap[event.(DataEvent).Name].Update(event.(DataEvent).Ohlc)

					//Run only once to setup indicators
					o.Do(func() {
						log.Debug("GENERATIONK>RUN ONCE")
						k.market.Strategy[0].Setup(k.market)
						log.WithFields(log.Fields{
							"strategy": k.market.Strategy[0],
						}).Debug("Strategy")
					})

					//Run setup after initperiod is finished
					if k.market.K < k.market.GetInitPeriod() {

						log.Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> Initializing strategy failed")
						break

					} else {

						log.Debug("GENERATIONK>EVENTCHANNEL> Updating indicators data")
						updateIndicators(k.market, event.(DataEvent))

						log.Debug("GENERATIONK>EVENTCHANNEL> Leting strategy know")
						k.market.Strategy[0].Tick(k)
					}
				case Quit:
					log.Debug("GENERATIONK>EVENTCHANNEL> QUIT EVENT PICKED OFF QUEUE")
					close(k.market.OrderChannel())
					close(k.market.EventChannel())
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
	//log.Debug("ctx.AssetIndicatorMap[dataEvent.Name]: ", len(ctx.AssetIndicatorMap[dataEvent.Name]))

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

		/*log.WithFields(log.Fields{
			"len(dataWindow)": len(dataWindow),
			"dataWindow":      dataWindow,
		}).Debug("GENERATIONK>UPDATE INDICATORS>")*/

		//Update the indicator with new data
		indicator.Update(dataWindow)
	}
}

type EndOfDataError struct {
	Description string
}

func (e *EndOfDataError) Error() string {
	return fmt.Sprintf("End of data: %s", e.Description)
}
