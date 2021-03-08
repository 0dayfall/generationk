package generationk

import (
	"errors"
	"fmt"
	"time"
)

//var o sync.Once

var AssetDoesNotExist = errors.New("Asset does not exist")
var FFToStartDate = errors.New("Fast forwarding to start date")
var EndOfBacktest = errors.New("End of backtest")
var EndOfData = errors.New("End of data")
var UnstablePeriod = errors.New("The stable period is not yet reached")
var Initialization = errors.New("Initialization in Once() failed")

//Callback is used in the strategy to give actions back to the backtest in progress
type Callback interface {
	Owning() bool
	IsOwning(assetName string) (bool, error)
	SendOrder(direction Direction, orderType OrderType, qty int) error
	SendOrderFor(assetName string, direction Direction, orderType OrderType, qty int) error
	Assets() []string
}

type GenerationK struct {
	ctx *Context
}

//NewGenerationK is used to create a new backtest
func NewGenerationK() *GenerationK {

	generationK := &GenerationK{
		ctx: NewContext(),
	}

	return generationK
}

func (g *GenerationK) inc() {
	g.ctx.K++
}

func (g *GenerationK) nextGen() error {
	defer g.inc()
	g.ctx.datePointer = g.ctx.asset.ohlc.Time[g.ctx.K]

	//Have to run this first so that we dont increase k by FF
	if g.ctx.K < 1 {
		//fmt.Printf("Once executed for %s\n\n", g.ctx.asset.name)
		err := g.ctx.strategy[0].Once(g.ctx, g.ctx.asset.ohlc)
		if err != nil {
			return err
		}
	}

	if g.ctx.asset.ohlc.Time[g.ctx.K].Before(g.ctx.startDate) {
		return FFToStartDate
	}

	if g.ctx.asset.ohlc.Time[g.ctx.K].After(g.ctx.endDate) {
		return EndOfBacktest
	}

	//Run setup after initperiod is finished
	if g.ctx.K < g.ctx.GetInitPeriod() {
		return UnstablePeriod
	}

	return g.ctx.strategy[0].PerBar(g.ctx.K, g)
}

func (k *GenerationK) Run() error {

	for k.ctx.K < k.ctx.length-1 {

		err := k.nextGen()
		if err != nil {
			switch err {
			case EndOfBacktest:
				return err
			case FFToStartDate:
				continue
			case UnstablePeriod:
				continue
			default:
				fmt.Println(err)
				return err
			}
		}

	}

	return nil
}

//AddDataManager is currently not used
func (k *GenerationK) SetDataManager() {}

//Returns an array of all assets
func (k *GenerationK) GetAsset() Asset {
	return k.ctx.GetAssets()[0]
}

//Returns an array of all assets
func (k *GenerationK) GetAssets() []Asset {
	return k.ctx.GetAssets()
}

//GetAssetByName returns a pointer to the asset by that name
func (k *GenerationK) GetAssetByName(name string) *Asset {
	return k.ctx.GetAssetByName(name)
}

//SetComission is used to set the comission scheme is there is one
func (k *GenerationK) SetComission(comission Comission) {
	k.ctx.broker.SetComission(comission)
}

//AddAsset is used to add a pointer to an asset
func (k *GenerationK) AddAsset(asset *Asset) {
	k.ctx.AddAsset(asset)
}

//AddPortfolio is used to add a pointer to a portfolio to the backtest
func (k *GenerationK) SetPortfolio(portfolio *Portfolio) {
	k.ctx.portfolio = portfolio
	k.ctx.broker.portfolio = portfolio
}

//AddStrategy is used to add a strategy to the backtest
func (k *GenerationK) AddStrategy(strat Strategy) {
	k.ctx.AddStrategy(strat)
}

//SetBalance is used to set the balance when the backtest is started
func (k *GenerationK) SetBalance(balance float64) {
	k.ctx.portfolio.SetBalance(balance)
}

//AddStartDate is used to set the end date for the backtest
func (k *GenerationK) SetStartDate(startDate time.Time) {
	k.ctx.SetStartDate(startDate)
}

//AddEndDate is used to set the end date for the backtest
func (k *GenerationK) SetEndDate(endDate time.Time) {
	k.ctx.SetEndDate(endDate)
}

//OrderSend is used to send an order to the broker, return an error if the asset does not exist
func (k *GenerationK) SendOrderFor(assetName string, direction Direction, orderType OrderType, qty int) error {
	if asset, ok := k.ctx.assetMap[k.ctx.asset.name]; ok {
		return k.sendOrder(k.ctx, direction, orderType, asset, k.ctx.datePointer, qty)
	}

	return AssetDoesNotExist
}

//OrderSend is used to send an order to the broker, return an error if the asset does not exist
func (k *GenerationK) SendOrder(direction Direction, orderType OrderType, qty int) error {
	if asset, ok := k.ctx.assetMap[k.ctx.asset.name]; ok {
		return k.sendOrder(k.ctx, direction, orderType, asset, k.ctx.datePointer, qty)
	}

	return AssetDoesNotExist
}

//orderSend is used to send an order to the broker
func (k *GenerationK) sendOrder(ctx *Context, direction Direction, orderType OrderType, asset *Asset, time time.Time, qty int) error {
	orderStatus, _ := interface{}(ctx.strategy[0]).(OrderStatus)

	err := ctx.broker.SendOrder(
		Order{
			direction: direction,
			orderType: orderType,
			Asset:     asset.name,
			Price:     asset.ohlc.Close[ctx.K],
			Time:      time,
			Qty:       qty,
		},
		orderStatus,
	)
	return err
}

//Assets returns an array of assets
func (k *GenerationK) Assets() []string {
	assets := make([]string, len(k.ctx.assets))
	for i, asset := range k.ctx.assets {
		assets[i] = asset.name
	}

	return assets
}

/*func (k *GenerationK) SetUniverse(assets []string) {
	for i := range assets {
		k.AddAsset(NewAsset(assets[i], nil))
	}
}*/

//OwnPosition is used to find out if we have a holding in an asset
//and the assumption is that the strategy is using multiple assets
func (k *GenerationK) IsOwning(assetName string) (bool, error) {
	if _, ok := k.ctx.assetMap[assetName]; ok {
		return k.ctx.portfolio.IsOwning(assetName), nil
	}
	return false, AssetDoesNotExist
}

//Owning is used to find out if we have a holding and we are
//only processing 1 asset
func (k *GenerationK) Owning() bool {
	return k.ctx.portfolio.IsOwning(k.ctx.asset.name)
}

//min returns the smaller of x or y.
func min(x, y int) int {
	if x > y {
		return y
	}

	return x
}
