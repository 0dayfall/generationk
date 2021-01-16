package generationk

import (
	"time"

	indicators "github.com/0dayfall/generationk/indicators"
)

//Context holds holds the strategy, assets, indicators per asset, start date, end date, the portfolio
//the current date and the unstable period
type Context struct {
	strategy []Strategy
	//	assetName         string
	asset             *Asset
	assets            []Asset
	assetMap          map[string]*Asset
	assetIndicatorMap map[string][]indicators.Indicator
	startDate         time.Time
	endDate           time.Time
	portfolio         *Portfolio
	broker            Broker
	K                 int
	datePointer       time.Time
	initPeriod        int
	length            int
}

//NewContext creates a new context
func NewContext() *Context {
	ctx := &Context{
		assets:            []Asset{},
		assetMap:          make(map[string]*Asset),
		assetIndicatorMap: make(map[string][]indicators.Indicator),
		broker:            Broker{},
		K:                 0,          //The first evolution
		initPeriod:        -1,         //Dont update strategy until K > initPeriod
		endDate:           time.Now(), //We are not time travellers, only back testers
		length:            0,
	}

	return ctx
}

//Time returns the time of the current bar
func (ctx *Context) Time() time.Time {
	return ctx.datePointer
}

//AddIndicatorOnAsset will add an indicator on the asset
func (ctx *Context) AddIndicatorOnAsset(asset *Asset, indicator indicators.Indicator) {
	ctx.assetIndicatorMap[asset.name] = append(ctx.assetIndicatorMap[asset.name], indicator)
}

//AddIndicator will add it to all assets
func (ctx *Context) AddIndicator(indicator indicators.Indicator) {
	for name := range ctx.assetMap {
		ctx.assetIndicatorMap[name] = append(ctx.assetIndicatorMap[name], indicator)
	}
}

//AddIndicator will add it to all assets
func (ctx *Context) AddIndicatorWithParams(indicator indicators.Indicator, param indicators.Param) {
	for name := range ctx.assetMap {
		ctx.assetIndicatorMap[name] = append(ctx.assetIndicatorMap[name], indicator)
	}
}

//AddEndDate is used to set the strategy that will be run
func (ctx *Context) AddEndDate(endTime time.Time) {
	ctx.endDate = endTime
}

//AddStartDate is used to set the start date
func (ctx *Context) AddStartDate(startTime time.Time) {
	ctx.startDate = startTime
	ctx.datePointer = startTime
}

//AddStrategy is used to set the strategy that will be run
func (ctx *Context) AddStrategy(strategy Strategy) {
	ctx.strategy = append(ctx.strategy, strategy)
}

//SetInitPeriod is used to set the unstable period, the longest period
//shoul be used
func (ctx *Context) SetInitPeriod(period int) {
	ctx.initPeriod = period
}

//GetInitPeriod returns the period
func (ctx *Context) GetInitPeriod() int {
	return ctx.initPeriod
}

//GetAssets returns the assets used in the strategy
func (ctx *Context) GetAssets() []Asset {
	return ctx.assets
}

//GetAssetByName return a specific strategy
func (ctx *Context) GetAssetByName(name string) *Asset {
	return ctx.assetMap[name]
}

//GetAssetIndicatorByName is used to get the indicators associated with the asset
func (ctx *Context) GetAssetIndicatorByName(name string) []indicators.Indicator {
	return ctx.assetIndicatorMap[name]
}

//AddAsset is used to add assets that the strategy will use
func (ctx *Context) AddAsset(asset *Asset) {
	//fmt.Printf("Adding asset: %s\n\n", asset.name)
	ctx.asset = asset
	ctx.assets = append(ctx.assets, *asset)
	ctx.assetMap[asset.name] = asset
	ctx.assetIndicatorMap[asset.name] = nil

	//Save the length of the longest asset
	if ctx.length < ctx.asset.length {
		ctx.length = ctx.asset.length
		//fmt.Printf("Length of asset after adding in ctx %d\n\n", ctx.length)
	}
	/*length := len(asset.ohlc.Close)
	if length > ctx.K {
		ctx.K = length
	}*/
}
