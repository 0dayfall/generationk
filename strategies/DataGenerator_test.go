package strategy

import (
	"generationk"
	"generationk/indicators"
	ind "generationk/indicators"
	genk "generationk/internal"
	"os"
	"testing"
	"time"

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
func (ma MACrossStrategy) Setup(ctx *genk.Context) error {
	ma.close = ind.NewTimeSeries(ind.Close, 5)
	ma.ma50 = ind.NewSimpleMovingAverage(ind.Close, 7)

	ctx.AddIndicator(ma.close)
	ctx.AddIndicator(ma.ma50)
	ctx.SetInitPeriod(7)

	return nil
}

//Update gets called when updates arrive
func (ma *MACrossStrategy) Update(ctx *genk.Context) {

	ctx.K++
}

//Tick get called when there is new data coming in
func (ma *MACrossStrategy) Tick(ctx *genk.Context) {

	if ma.ma50.ValueAtIndex(0) > ma.close.ValueAtIndex(0) {
		if !ctx.Position(ctx.AssetMap["ABB"]) {
			generationk.MakeOrder(ctx, genk.OrderType(genk.Buy), ctx.AssetMap["ABB"], ctx.Time(), 1000)
		}
	}

	if ma.ma50.ValueAtIndex(0) < ma.close.ValueAtIndex(0) {
		MakeOrder(ctx, genk.OrderType(genk.Sell), ctx.AssetMap["ABB"], ctx.Time(), 1000)
	}

}

//Orders get called when everything is updated
func (ma *MACrossStrategy) OrderEvent(ctx *genk.Context) {
	log.Debug("MAStrategy_test> OrderEvent")
}

func TestRun(t *testing.T) {

	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	//logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	lvl, ok := os.LookupEnv("LOG_LEVEL")

	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	// set global log level
	logrus.SetLevel(ll)

	//Context that the strategy is being run with such as assets
	market := genk.NewContext()
	market.AddAsset(genk.NewAsset("ABB", genk.OHLC{}))
	//Going to run with the following data thingie to collect the data
	dataManager := genk.NewCSVDataManager(market)
	dataManager.ReadCSVFileAsync("test/data/ABB.csv")
	strategy := genk.Strategy(&MACrossStrategy{})
	market.AddStrategy(&strategy)

	now := time.Now()
	start := now.AddDate(0, -9, -2)
	market.AddStartDate(start)

	now = time.Now()
	end := now.AddDate(0, -3, -2)
	market.AddStartDate(end)

	generationk.RunEventBased(market)
}
