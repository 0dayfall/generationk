package generationk

import (
	"fmt"
	"generationk/indicators"
	ind "generationk/indicators"
	genk "generationk/internal"
	"os"
	"testing"
	"time"

	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

//Strategy strategy
type MACrossStrategy struct {
	ma50  *indicators.Average
	close *indicators.Series
}

//Setup is used to start the strategy
func (m *MACrossStrategy) Setup(ctx *genk.Context) {
	fmt.Printf("Init strategy\n")
	m.close = ind.TimeSeries(ctx.AssetMap["ABB"].CloseArray())
	m.ma50 = ind.SimpleMovingAverage(ctx.AssetMap["ABB"].CloseArray(), 50)
	//ma200 := *ind.SimpleMovingAverage(ctx.AssetMap["ABB"].CloseArray(), 200)
}

//Update gets called when updates arrive
func (m *MACrossStrategy) Update(ctx *genk.Context) {
	//New day new values
	ctx.K++
}

//Orders get called when everything is updated
func (m *MACrossStrategy) Tick(ctx *genk.Context) {
	if m.ma50.ValueAtIndex(ctx.K) > m.close.ValueAtIndex(ctx.K) {
		MakeOrder(ctx, genk.OrderType(genk.Buy), ctx.AssetMap["ABB"], ctx.Time(), 1000)
	}
	ctx.K++
	//generationk.Signal()
}

//Orders get called when everything is updated
func (m *MACrossStrategy) OrderEvent(ctx *genk.Context) {
	log.Debug("MAStrategy_test> OrderEvent")
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
	market := genk.NewContext()

	//Going to run with the following data
	dataManager := genk.NewDataManager(market.EventChannel())
	

	strategy := genk.Strategy(&MACrossStrategy{})
	market.AddStrategy(&strategy)

	now := time.Now()
	then := now.AddDate(0, -9, -2)
	market.AddStartDate(then)

	go RunBacktest(market)
	select {}
}
