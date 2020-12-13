package generationk

import (
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

//Setup is used to start the strategy
func (m *MACrossStrategy) Setup(ctx *genk.Context) error {
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
	m.close = ind.NewTimeSeries(ind.Close, 5)
	m.ma50 = ind.NewSimpleMovingAverage(ind.Close, 7)

	ctx.AddIndicator(m.close)
	ctx.AddIndicator(m.ma50)
	m.SetInitPeriod(7)
	//ma200 := *ind.SimpleMovingAverage(ctx.AssetMap["ABB"].CloseArray(), 200)
	return nil
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
		if !ctx.Position(ctx.AssetMap["ABB"]) {
			MakeOrder(ctx, genk.OrderType(genk.Buy), ctx.AssetMap["ABB"], ctx.Time(), 1000)
		}
	}
	if m.ma50.ValueAtIndex(0) < m.close.ValueAtIndex(0) {
		MakeOrder(ctx, genk.OrderType(genk.Sell), ctx.AssetMap["ABB"], ctx.Time(), 1000)
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
	go dataManager.ReadCSVFile("test/data/ABB.csv")
	strategy := genk.Strategy(&MACrossStrategy{})
	market.AddStrategy(&strategy)

	now := time.Now()
	start := now.AddDate(0, -9, -2)
	market.AddStartDate(start)

	now = time.Now()
	end := now.AddDate(0, -3, -2)
	market.AddStartDate(end)

	RunLive(market)
}
