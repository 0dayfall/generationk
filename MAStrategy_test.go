package generationk

import (
	"fmt"
	genk "generationk/internal"
	"testing"
)

//Strategy strategy
type MACrossStrategy struct {
}

//Init is used to start the strategy
func (m *MACrossStrategy) Init(ctx *genk.Context) {
	fmt.Printf("Init strategy\n")
	//ma := indicators.SimpleMovingAverage(ctx, 9)
}

//Tick gets called when new data is available
func (m *MACrossStrategy) Tick(ctx *genk.Context) {
	fmt.Printf("Close: %f\n", ctx.Asset[0].Ohlc[0].Close)
}

func TestRun(t *testing.T) {
	dataManager := genk.NewDataManager()
	var strategy genk.Strategy
	var market genk.Context
	var asset genk.Asset

	strategy = &MACrossStrategy{}
	asset = dataManager.ReadCSVFile("test/ABB.csv")
	market.AddAsset(&asset)
	market.AddStrategy(&strategy)

	go RunBacktest(&market)
	select {}
}
