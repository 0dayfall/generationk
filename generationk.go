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

var AssetDoesNotExist = errors.New("Asset does not exist")

type GenkCallback interface {
	IsOwning(asset string) bool
	OrderSend(assetName string, direction Directon, orderType OrderType, amount float64, qty int) error
	Assets() []string
}

type GenerationK struct {
	market *Context
}

func NewGenerationK() *GenerationK {

	generationK := &GenerationK{
		market: NewContext(),
	}

	generationK.initLog()

	return generationK
}

func (k *GenerationK) UpdateIndicators(assetName string) {
	k.market.updateIndicators(assetName)
}

func (k *GenerationK) DataEvent(dataEvent Event) {
	assetName := dataEvent.(DataEvent).Name
	/*log.WithFields(log.Fields{
		"Number of items": len(k.market.EventChannel()),
	}).Debug("GENERATIONK>DATA EVENT PICKED OFF QUEUE")*/

	k.market.datePointer = dataEvent.(DataEvent).Ohlc.time

	if dataEvent.(DataEvent).Ohlc.time.After(k.market.endDate) {
		log.Debug("GENERATIONK>EVENTCHANNEL> Ohlc.Time is after the back test end date")
		//k.market.EventChannel() <- Quit{}
		return
	}

	//Add data to asset
	if k.market.assetMap == nil {
		log.Info("ASSET MAP == NIL")
	}

	if _, ok := k.market.assetMap[assetName]; !ok {
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
		k.GetAssetByName(assetName).Update(dataEvent.(DataEvent).Ohlc, k.market.initPeriod)
		k.UpdateIndicators(dataEvent.(DataEvent).Name)
	}
	//Run only once to setup indicators
	/*o.Do(func() {
	})*/

	if k.market.K < 1 {

		log.Info("GENERATIONK>RUN ONCE")

		err := k.market.strategy[0].Once(k.market)
		if err != nil {
			log.Fatal(err)
		}

		//Check if start date is set, otherwise use the date of the first asset ohlc bar
		//Check if the end date is set, otherwise use the data of today

		log.WithFields(log.Fields{
			"strategy": k.market.strategy[0],
		}).Debug("Strategy")
	}
	k.market.K++

	//Run setup after initperiod is finished
	if k.market.K < k.market.GetInitPeriod() {
		log.Debug("GENERATIONK>EVENTCHANNEL>DATAEVENT> Strategy in init period")

		return
	}

	log.Debug("GENERATIONK>EVENTCHANNEL> Updating indicators data")
	log.Debug("GENERATIONK>EVENTCHANNEL> Leting strategy know")
	k.market.strategy[0].PerBar(k)

	log.WithFields(log.Fields{
		"K: ": k.market.K,
	}).Debug("K")

}

func (k *GenerationK) initLog() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	lvl, ok := os.LookupEnv("LOG_LEVEL")

	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "fatal"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.FatalLevel
	}
	// set global log level
	logrus.SetLevel(ll)
}

func (k *GenerationK) AddDataManager() {}

func (k *GenerationK) GetAssets() []Asset {
	return k.market.GetAssets()
}

func (k *GenerationK) GetAssetByName(name string) *Asset {
	return k.market.GetAssetByName(name)
}

func (k *GenerationK) AddComission(comission Comission) {
	k.market.broker.SetComission(comission)
}

func (k *GenerationK) AddAsset(asset *Asset) {
	k.market.AddAsset(asset)
}

func (k *GenerationK) AddPortfolio(portfolio *Portfolio) {
	k.market.portfolio = portfolio
	k.market.broker.portfolio = portfolio
}

func (k *GenerationK) AddStrategy(strat Strategy) {

	/*err := strat.Setup(k.market)
	if err != nil {
		log.Fatal("Could not initialize strategy")
	}*/
	k.market.AddStrategy(strat)
}

func (k *GenerationK) SetBalance(balance float64) {
	k.market.portfolio.SetBalance(balance)
}

func (k *GenerationK) AddStartDate(startDate time.Time) {
	k.market.AddStartDate(startDate)
}

func (k *GenerationK) AddEndDate(endDate time.Time) {
	k.market.AddEndDate(endDate)
}

func (k *GenerationK) OrderSend(assetName string, direction Directon, orderType OrderType, amount float64, qty int) error {
	if asset, ok := k.market.assetMap[assetName]; ok {
		//do something here
		orderSend(k.market, direction, orderType, asset, k.market.datePointer, amount, qty)

		return nil
	}
	return AssetDoesNotExist
}

func orderSend(ctx *Context, direction Directon, orderType OrderType, asset *Asset, time time.Time, amount float64, qty int) {
	/*log.WithFields(log.Fields{
		"Direction":  direction,
		"Order type": orderType,
		"Time":       time,
		"Amount":     amount,
		"Qty":        qty,
	}).Info("GENERATIONK>MAKE ORDER>")*/

	orderStatus, _ := interface{}(ctx.strategy[0]).(OrderStatus)

	ctx.broker.SendOrder(
		Order{
			direction: direction,
			orderType: orderType,
			Asset:     asset,
			Time:      time,
			Amount:    amount,
			Qty:       qty,
		},
		orderStatus,
	)

}

func (k *GenerationK) Assets() []string {
	log.WithFields(log.Fields{
		"Length": len(k.market.assets),
	}).Debug("Length of assets field")

	assets := make([]string, len(k.market.assets))
	for i, asset := range k.market.assets {
		assets[i] = asset.name
	}

	return assets
	//return []string{"Test", "Test2"}
}

//OwnPosition is used to find out if we have a holding in an asset
func (k *GenerationK) IsOwning(name string) bool {
	return k.market.portfolio.IsOwning(name)
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
