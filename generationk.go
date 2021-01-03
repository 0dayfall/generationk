package generationk

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

var o sync.Once

var AssetDoesNotExist = errors.New("Asset does not exist")

//Callback is used in the strategy to give actions back to the backtest in progress
type Callback interface {
	IsOwning(asset string) bool
	OrderSend(assetName string, direction Direction, orderType OrderType, amount float64, qty int) error
	Assets() []string
}

type GenerationK struct {
	market *Context
}

//NewGenerationK is used to create a new backtest
func NewGenerationK() *GenerationK {

	generationK := &GenerationK{
		market: NewContext(),
	}

	return generationK
}

//Is used to update the indicator with the assetName
func (k *GenerationK) UpdateIndicators(assetName string) {
	k.market.updateIndicators(assetName)
}

//DataEvent is used to implement the datahandler interface and called when reading a file
func (k *GenerationK) DataEvent(dataEvent Event) {
	assetName := dataEvent.(DataEvent).Name
	k.market.datePointer = dataEvent.(DataEvent).Ohlc.Time

	if dataEvent.(DataEvent).Ohlc.Time.After(k.market.endDate) || dataEvent.(DataEvent).Ohlc.Time.Before(k.market.startDate) {
		return
	}

	if _, ok := k.market.assetMap[assetName]; !ok {
		asset := NewAsset(assetName, dataEvent.(DataEvent).Ohlc)
		k.AddAsset(asset)
	} else {
		k.GetAssetByName(assetName).Update(dataEvent.(DataEvent).Ohlc, k.market.initPeriod)
		k.UpdateIndicators(dataEvent.(DataEvent).Name)
	}

	if k.market.K < 1 {
		err := k.market.strategy[0].Once(k.market)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	k.market.K++

	//Run setup after initperiod is finished
	if k.market.K < k.market.GetInitPeriod() {
		return
	}

	k.market.strategy[0].PerBar(dataEvent.(DataEvent).Ohlc, (k))
}

//AddDataManager is currently not used
func (k *GenerationK) AddDataManager() {}

//Returns an array of all assets
func (k *GenerationK) GetAssets() []Asset {
	return k.market.GetAssets()
}

//GetAssetByName returns a pointer to the asset by that name
func (k *GenerationK) GetAssetByName(name string) *Asset {
	return k.market.GetAssetByName(name)
}

//SetComission is used to set the comission scheme is there is one
func (k *GenerationK) AddComission(comission Comission) {
	k.market.broker.SetComission(comission)
}

//AddAsset is used to add a pointer to an asset
func (k *GenerationK) AddAsset(asset *Asset) {
	k.market.AddAsset(asset)
}

//AddPortfolio is used to add a pointer to a portfolio to the backtest
func (k *GenerationK) AddPortfolio(portfolio *Portfolio) {
	k.market.portfolio = portfolio
	k.market.broker.portfolio = portfolio
}

//AddStrategy is used to add a strategy to the backtest
func (k *GenerationK) AddStrategy(strat Strategy) {
	k.market.AddStrategy(strat)
}

//SetBalance is used to set the balance when the backtest is started
func (k *GenerationK) SetBalance(balance float64) {
	k.market.portfolio.SetBalance(balance)
}

//AddStartDate is used to set the end date for the backtest
func (k *GenerationK) AddStartDate(startDate time.Time) {
	k.market.AddStartDate(startDate)
}

//AddEndDate is used to set the end date for the backtest
func (k *GenerationK) AddEndDate(endDate time.Time) {
	k.market.AddEndDate(endDate)
}

//OrderSend is used to send an order to the broker, return an error if the asset does not exist
func (k *GenerationK) OrderSend(assetName string, direction Direction, orderType OrderType, amount float64, qty int) error {
	if asset, ok := k.market.assetMap[assetName]; ok {
		orderSend(k.market, direction, orderType, asset, k.market.datePointer, amount, qty)

		return nil
	}
	return AssetDoesNotExist
}

//orderSend is used to send an order to the broker
func orderSend(ctx *Context, direction Direction, orderType OrderType, asset *Asset, time time.Time, amount float64, qty int) {
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

//Assets returns an array of assets
func (k *GenerationK) Assets() []string {
	assets := make([]string, len(k.market.assets))
	for i, asset := range k.market.assets {
		assets[i] = asset.name
	}

	return assets
}

//OwnPosition is used to find out if we have a holding in an asset
//and the assumption is that the strategy is using multiple assets
func (k *GenerationK) IsOwning(name string) bool {
	return k.market.portfolio.IsOwning(name)
}

//min returns the smaller of x or y.
func min(x, y int) int {
	if x > y {
		return y
	}

	return x
}

//Owning is used to find out if we have a holding and we are
//only processing 1 asset
func (k *GenerationK) Owning() bool {
	return k.market.portfolio.IsOwning(k.market.assets[0].name)
}

type EndOfDataError struct {
	Description string
}

func (e *EndOfDataError) Error() string {
	return fmt.Sprintf("End of data: %s", e.Description)
}

func RunStrategyOnAssets(strategy Strategy, folderPath string) {
	files, err := filepath.Glob(folderPath + "*.csv")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("files %s", files)
	//d.ReadCSVFilesAsync(files)
	portfolio := NewPortfolio()
	portfolio.SetBalance(100000)

	var wg sync.WaitGroup

	y := 0

	for _, fileName := range files {
		wg.Add(1)
		go func(localFilename string) {
			genk := NewGenerationK()
			genk.AddPortfolio(portfolio)
			genk.AddStrategy(strategy)

			now := time.Now()
			start := now.AddDate(-15, -9, -2)
			genk.AddStartDate(start)
			now = time.Now()
			end := now.AddDate(0, -3, -2)
			genk.AddEndDate(end)

			//genk.RunEventBased()
			dataManager := NewCSVDataManager(genk)
			//dataManager.SetHandler(genk)
			//genk.AddDataManager(dataManager)

			//dataManager.ReadCSVFilesAsync([]string{"test/data/ABB.csv", "test/data/ASSAb.csv"})
			dataManager.ReadCSVFile(localFilename)

			wg.Done()
		}(fileName)
		y++
	}
	wg.Wait()
}
