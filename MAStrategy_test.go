package generationk

import (
	"testing"
	"time"

	indicators "github.com/0dayfall/generationk/indicators"

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
	//Want access to the latest 5 closing prices
	ma.close = indicators.NewTimeSeries(indicators.Close, 5)
	//MA50
	ma.ma50 = indicators.NewSimpleMovingAverage(indicators.Close, 50)

	//Add indicators to context
	ctx.AddIndicator(ma.close)
	ctx.AddIndicator(ma.ma50)

	//The data needed to calculate MA
	ctx.SetInitPeriod(50)

	return nil
}

//Update gets called when updates arrive
func (ma *MACrossStrategy) Update(ctx *Context) {
	ctx.K++
}

//Tick get called when there is new data coming in
func (ma *MACrossStrategy) Tick(broker GenkCallback) {

	if ma.close.ValueAtIndex(0) > ma.ma50.ValueAtIndex(0) {
		if !broker.IsOwning("ABB") {
			broker.OrderSend("ABB", OrderType(BuyOrder), 0, 100)
		}
	}

	if ma.close.ValueAtIndex(0) < ma.ma50.ValueAtIndex(0) {
		if broker.IsOwning("ABB") {
			broker.OrderSend("ABB", OrderType(SellOrder), 0, 100)
		}
	}

}

//OrderEvent gets called on order events
func (ma *MACrossStrategy) OrderEvent(orderEvent Event) {
	log.WithFields(log.Fields{
		"orderEvent": orderEvent,
	}).Debug("MAStrategy_test> OrderEvent")
}

func TestRun(t *testing.T) {

	genk := NewGenerationK()
	genk.Init()

	strategy := Strategy(&MACrossStrategy{})
	//Going to run with the following data thingie to collect the data
	genk.AddAsset(NewAsset("ABB", OHLC{}))
	genk.AddAsset(NewAsset("ASSAb", OHLC{}))

	genk.AddStrategy(&strategy)
	genk.SetBalance(100000)

	now := time.Now()
	start := now.AddDate(0, -9, -2)
	genk.AddStartDate(start)

	now = time.Now()
	end := now.AddDate(0, -3, -2)
	genk.AddEndDate(end)

	//genk.RunEventBased()
	dataManager := NewCSVDataManager(genk)
	dataManager.ReadCSVFilesAsync([]string{"test/data/ABB.csv", "test/data/ASSAb.csv"})
}

func BenchmarkRun(t *testing.B) {

	genk := NewGenerationK()
	genk.Init()

	strategy := Strategy(&MACrossStrategy{})
	//Going to run with the following data thingie to collect the data
	genk.AddAsset(NewAsset("ABB", OHLC{}))

	genk.AddStrategy(&strategy)
	genk.SetBalance(100000)

	now := time.Now()
	start := now.AddDate(0, -9, -2)
	genk.AddStartDate(start)

	now = time.Now()
	end := now.AddDate(0, -3, -2)
	genk.AddEndDate(end)

	//genk.RunEventBased()
	dataManager := NewCSVDataManager(genk)
	dataManager.ReadCSVFile("test/data/ABB.csv")
}
