package generationk

import (
	"fmt"
	indicators "generationk/indicators"
	genk "generationk/internal"
	"testing"
	"time"
)

//Strategy strategy
type MACrossStrategy struct {
	ma50 indicators.Average
}

//Init is used to start the strategy
func (m *MACrossStrategy) Init(ctx *genk.Context) {
	fmt.Printf("Init strategy\n")
	m.ma50 = *indicators.SimpleMovingAverage(ctx.AssetMap["ABB"].CloseArray(), 50)
	ctx.A
}

//Tick gets called when new data is available
func (m *MACrossStrategy) Tick(ctx *genk.Context) {
	fmt.Printf("Date: %v ", ctx.Time())
	fmt.Printf("comparing %f with %f\n", m.ma50.Value(), ctx.AssetMap["ABB"].Close())
	fmt.Printf("len asset map %d \n", len(ctx.AssetMap["ABB"].Ohlc))
	if m.ma50.Value() > ctx.AssetMap["ABB"].Close() {
		fmt.Printf("BUY! ===============>")
	}
	//fmt.Printf("Close: %f\n", ctx.Asset[0].Ohlc[0].Close)
}

func TestRun(t *testing.T) {
	dataManager := genk.NewDataManager()
	market := genk.NewContext()
	strategy := genk.Strategy(&MACrossStrategy{})
	asset := dataManager.ReadCSVFile("test/ABB.csv")

	market.AddAsset(&asset)
	market.AddStrategy(&strategy)
	now := time.Now()
	then := now.AddDate(0, -5, -2)
	market.AddStartDate(then)

	go RunBacktest(market)
	select {}
}
