package generationk

import (
	"os"
	"testing"
	"time"

	indicators "github.com/greenorangebay/generationk/indicators"
	"github.com/shiena/ansicolor"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

//Strategy strategy
type MACrossStrategy struct {
	ma50       *indicators.SimpleMovingAverage
	close      *indicators.TimeSeries
	initPeriod int
}

//Setup is used to declare what indicators will be used
func (ma *MACrossStrategy) Setup(ctx *Context) error {
	ma.close = indicators.NewTimeSeries(indicators.Close, 5)
	ma.ma50 = indicators.NewSimpleMovingAverage(indicators.Close, 9)

	ctx.AddIndicator(ma.close)
	ctx.AddIndicator(ma.ma50)
	ctx.SetInitPeriod(9)

	return nil
}

//Update gets called when updates arrive
func (ma *MACrossStrategy) Update(ctx *Context) {
	ctx.K++
}

//Tick get called when there is new data coming in
func (ma *MACrossStrategy) Tick(ctx *Context) {

	if ma.ma50.ValueAtIndex(0) > ma.close.ValueAtIndex(0) {
		if !ctx.Position(ctx.AssetMap["ABB"]) {
			MakeOrder(ctx, OrderType(Buy), ctx.AssetMap["ABB"], ctx.Time(), 1000)
		}
	}

}

//Orders get called when everything is updated
func (ma *MACrossStrategy) OrderEvent(orderEvent Event) {
	log.WithFields(log.Fields{
		"orderEvent": orderEvent,
	}).Debug("MAStrategy_test> OrderEvent")
}

func TestRun(t *testing.T) {

	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	lvl, ok := os.LookupEnv("LOG_LEVEL")

	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "debug"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	// set global log level
	logrus.SetLevel(ll)

	//Context that the strategy is being run with such as assets
	market := NewContext()
	market.AddAsset(NewAsset("ABB", OHLC{}))
	//Going to run with the following data thingie to collect the data
	dataManager := NewCSVDataManager(market)
	dataManager.ReadCSVFileAsync("test/data/ABB.csv")
	strategy := Strategy(&MACrossStrategy{})
	market.AddStrategy(&strategy)

	now := time.Now()
	start := now.AddDate(0, -9, -2)
	market.AddStartDate(start)

	now = time.Now()
	end := now.AddDate(0, -3, -2)
	market.AddStartDate(end)

	RunEventBased(market)
}
