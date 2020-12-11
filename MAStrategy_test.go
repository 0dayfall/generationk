package generationk

import (
	"generationk/indicators"
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
	close      *indicators.Series
	initPeriod int
}

//Setup is used to start the strategy
func (m *MACrossStrategy) Setup(ctx *genk.Context) {
	/*var e error
	if ctx.K < 5 {
		return errors.New("Need more data to calculate indicators")
	}
	m.close, e = ind.NewTimeSeries(ctx.AssetMap["ABB"])
	m.ma50, e = ind.NewSimpleMovingAverage(ctx.AssetMap["ABB"], ind.Close, 5)
	if e != nil {
		return e
	}
	ctx.AddUpdatable(m.close, m.ma50)

	return nil*/
	m.close, e = ind.NewTimeSeries(Close)
	m.ma50, e = ind.NewSimpleMovingAverage(genk.Close, 5)

	ctx.AddIndicator(m.close)
	ctx.AddIndicator(m.m50)
	//ma200 := *ind.SimpleMovingAverage(ctx.AssetMap["ABB"].CloseArray(), 200)
}

func (m *MACrossStrategy) SetInitPeriod(period int) {
	m.initPeriod = period
}

func (m MACrossStrategy) GetInitPeriod() int {
	return m.initPeriod
}

//Update gets called when updates arrive
func (m *MACrossStrategy) Update(ctx *genk.Context) {
	//New day new values
	log.Debug("MASTRATEGY_TEST> Update")
	ctx.K++
	//m.close, _ = ind.TimeSeries(ctx.AssetMap["ABB"].CloseArray())
	//m.ma50, _ = ind.SimpleMovingAverage(ctx.AssetMap["ABB"].CloseArray(), 5)
}

//Orders get called when everything is updated
func (m *MACrossStrategy) Tick(ctx *genk.Context) {
	if m.ma50.ValueAtIndex(0) > m.close.ValueAtIndex(0) {
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

	//logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	//logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

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

	//Going to run with the following data thingie to collect the data
	dataManager := genk.NewCSVDataManager(market)
	go dataManager.ReadCSVFile("test/data/ABB.csv")
	strategy := genk.Strategy(&MACrossStrategy{})
	market.AddStrategy(&strategy)

	now := time.Now()
	then := now.AddDate(0, -9, -2)
	market.AddStartDate(then)

	RunLive(market)
}
