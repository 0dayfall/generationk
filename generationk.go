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

	generationK.initLog()

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
	})
	if k.market.K < 1 {
		log.Info("GENERATIONK>RUN ONCE")
		k.market.Strategy[0].Setup(k.market)
		log.WithFields(log.Fields{
			"strategy": k.market.Strategy[0],
		}).Debug("Strategy")

	}
	k.market.K++

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

func (k *generationK) initLog() {
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

func (k *generationK) AddDataManager() {}

func (k *generationK) AddAsset(asset *Asset) {
	k.market.AddAsset(asset)
}

func (k *generationK) AddPortfolio(portfolio *Portfolio) {
	k.market.Portfolio = portfolio
	k.market.Broker.portfolio = portfolio
}

func (k *generationK) AddStrategy(strat *Strategy) {
	err := (*strat).Setup(k.market)
	if err != nil {
		log.Fatal("Could not initialize strategy")
	}
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
	}
}

type EndOfDataError struct {
	Description string
}

func (e *EndOfDataError) Error() string {
	return fmt.Sprintf("End of data: %s", e.Description)
}
