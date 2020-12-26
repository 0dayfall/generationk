package generationk

import (
	"errors"
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
	OrderSend(assetName string, ordertype OrderType, amount float64, qty int) error
	Assets() []string
}

type generationK struct {
	market *Context
}

func NewGenerationK() *generationK {

	generationK := &generationK{
		market: NewContext(),
	}

	generationK.initLog()

	return generationK
}

func (k *generationK) UpdateIndicators(assetName string) {
	k.market.updateIndicators(assetName)
}

func (k *generationK) DataEvent(dataEvent Event) {
	assetName := dataEvent.(DataEvent).Name
	/*log.WithFields(log.Fields{
		"Number of items": len(k.market.EventChannel()),
	}).Debug("GENERATIONK>DATA EVENT PICKED OFF QUEUE")*/

	k.market.datePointer = dataEvent.(DataEvent).Ohlc.Time

	if dataEvent.(DataEvent).Ohlc.Time.After(k.market.EndDate) {
		log.Debug("GENERATIONK>EVENTCHANNEL> Ohlc.Time is after the back test end date")
		//k.market.EventChannel() <- Quit{}
		return
	}

	//Add data to asset
	if k.market.AssetMap == nil {
		log.Info("ASSET MAP == NIL")
	}
	if _, ok := k.market.AssetMap[assetName]; !ok {
		log.Info("GENERATIONK>EVENTCHANNEL>DATAEVENT> CREATING ASSET AND ADDING TO MAP")
		asset := NewAsset(assetName, dataEvent.(DataEvent).Ohlc)
		k.AddAsset(asset)
		//k.market.Asset = append(k.market.Asset, *asset)
		//k.market.AssetMap[assetName] = asset
	} else {

		/*log.WithFields(log.Fields{
			"k.market.AssetMap[dataEvent.(DataEvent).Name].Name": k.market.AssetMap[dataEvent.(DataEvent).Name].Name,
			"k.market.AssetMap[dataEvent.(DataEvent).Name].Ohlc": k.market.AssetMap[dataEvent.(DataEvent).Name].Ohlc,
		}).Info("GENERATIONK>EVENTCHANNEL>DATAEVENT> EXISTS IN MAP")*/

		//k.market.AssetMap[event.(DataEvent).Name].Ohlc = prepend(k.market.AssetMap[event.(DataEvent).Name].Ohlc, event.(DataEvent).Ohlc)

		//k.market.AssetMap[assetName].Update(dataEvent.(DataEvent).Ohlc)
		k.GetAssetByName(assetName).Update(dataEvent.(DataEvent).Ohlc)
		k.UpdateIndicators(dataEvent.(DataEvent).Name)
	}
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
		log.Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> Strategy in init period")
		return
	} else {
		log.Debug("GENERATIONK>EVENTCHANNEL> Updating indicators data")
		log.Debug("GENERATIONK>EVENTCHANNEL> Leting strategy know")
		k.market.Strategy[0].Tick(k)
	}

	log.WithFields(log.Fields{
		"K: ": k.market.K,
	}).Debug("K")

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

func (k *generationK) GetAssets() []Asset {
	return k.market.GetAssets()
}

func (k *generationK) GetAssetByName(name string) *Asset {
	return k.market.GetAssetByName(name)
}

func (k *generationK) AddAsset(asset *Asset) {
	k.market.AddAsset(asset)
}

func (k *generationK) AddPortfolio(portfolio *Portfolio) {
	k.market.Portfolio = portfolio
	k.market.Broker.portfolio = portfolio
}

func (k *generationK) AddStrategy(strat Strategy) {
	/*err := strat.Setup(k.market)
	if err != nil {
		log.Fatal("Could not initialize strategy")
	}*/
	k.market.AddStrategy(strat)
}

func (k *generationK) SetBalance(balance float64) {
	k.market.Portfolio.SetBalance(balance)
}

func (k *generationK) AddStartDate(startDate time.Time) {
	k.market.AddStartDate(startDate)
}

func (k *generationK) AddEndDate(endDate time.Time) {
	k.market.AddEndDate(endDate)
}

func (k *generationK) OrderSend(assetName string, ordertype OrderType, amount float64, qty int) error {
	if asset, ok := k.market.AssetMap[assetName]; ok {
		//do something here
		orderSend(k.market, ordertype, asset, k.market.datePointer, amount, qty)
		return nil
	}
	return errors.New("Asset not in map")
}

func orderSend(ctx *Context, ordertype OrderType, asset *Asset, time time.Time, amount float64, qty int) {
	log.WithFields(log.Fields{
		"Order type": ordertype,
		"Time":       time,
		"Amount":     amount,
		"Qty":        qty,
	}).Info("GENERATIONK>MAKE ORDER>")

	orderStatus, _ := interface{}(ctx.Strategy[0]).(OrderStatus)

	ctx.Broker.SendOrder(
		Order{
			Ordertype: ordertype,
			Asset:     asset,
			Time:      time,
			Amount:    amount,
			Qty:       qty,
		},
		orderStatus,
	)

}

func (k *generationK) Assets() []string {
	log.WithFields(log.Fields{
		"Length": len(k.market.Asset),
	}).Debug("Length of assets field")
	assets := make([]string, len(k.market.Asset))
	for i, asset := range k.market.Asset {
		assets[i] = asset.Name
	}
	return assets
	//return []string{"Test", "Test2"}
}

//OwnPosition is used to find out if we have a holding in an asset
func (k *generationK) IsOwning(name string) bool {
	return k.market.Portfolio.IsOwning(name)
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

type EndOfDataError struct {
	Description string
}

func (e *EndOfDataError) Error() string {
	return fmt.Sprintf("End of data: %s", e.Description)
}
