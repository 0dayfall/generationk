package generationk

import (
	"time"

	D "github.com/0dayfall/generationk/data"
	I "github.com/0dayfall/generationk/indicators"
)

//Context holds holds the strategy, assets, indicators per asset, start date, end date, the portfolio
//the current date and the unstable period
type Context struct {
	strategy []Strategy
	//	assetName         string
	asset             *D.Asset
	assets            []*D.Asset
	assetMap          map[string]*D.Asset
	assetIndicatorMap map[string][]I.Indicator
	startDate         time.Time
	endDate           time.Time
	portfolio         *Portfolio
	broker            Broker
	K                 int
	datePointer       time.Time
	initPeriod        int
	length            int
	dataPath          string
}

//NewContext creates a new context
func NewContext() *Context {
	ctx := &Context{
		assets:            []*D.Asset{},
		assetMap:          make(map[string]*D.Asset),
		assetIndicatorMap: make(map[string][]I.Indicator),
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
func (ctx *Context) AddIndicatorOnAsset(asset *D.Asset, indicator I.Indicator) {
	ctx.assetIndicatorMap[asset.Name] = append(ctx.assetIndicatorMap[asset.Name], indicator)
}

//AddIndicator will add it to all assets
func (ctx *Context) AddIndicator(indicator I.Indicator) {
	for name := range ctx.assetMap {
		ctx.assetIndicatorMap[name] = append(ctx.assetIndicatorMap[name], indicator)
	}
}

//AddIndicator will add it to all assets
func (ctx *Context) AddIndicatorWithParams(indicator I.Indicator, param I.Param) {
	for name := range ctx.assetMap {
		ctx.assetIndicatorMap[name] = append(ctx.assetIndicatorMap[name], indicator)
	}
}

//AddEndDate is used to set the strategy that will be run
func (ctx *Context) SetEndDate(endTime time.Time) {
	ctx.endDate = endTime
}

func (ctx *Context) GetEndDate() time.Time {
	return ctx.endDate
}

//AddStartDate is used to set the start date
func (ctx *Context) SetStartDate(startTime time.Time) {
	ctx.startDate = startTime
	ctx.datePointer = startTime
}

func (ctx *Context) GetStartDate() time.Time {
	return ctx.startDate
}

//AddStrategy is used to set the strategy that will be run
func (ctx *Context) SetStrategy(strategy Strategy) {
	//Should be possible to use multiple strategies in the future
	ctx.strategy = append(ctx.strategy, strategy)
}

//SetStrategy is used to set the strategy that will be run
func (ctx *Context) GetStrategy() Strategy {
	// In the futures should return an array of
	// strategies when multiple strategies will be
	// handled
	return ctx.strategy[0]
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
func (ctx *Context) GetAssets() []*D.Asset {
	return ctx.assets
}

//GetAssetByName return a specific strategy
func (ctx *Context) GetAssetByName(name string) *D.Asset {
	return ctx.assetMap[name]
}

//GetAssetIndicatorByName is used to get the indicators associated with the asset
func (ctx *Context) GetAssetIndicatorByName(name string) []I.Indicator {
	return ctx.assetIndicatorMap[name]
}

//AddAsset is used to add assets that the strategy will use
func (ctx *Context) AddAsset(asset *D.Asset) {
	ctx.assets = append(ctx.assets, asset)
	ctx.assetMap[asset.Name] = asset
	ctx.assetIndicatorMap[asset.Name] = nil

	//Remove is below is commented back
	ctx.length = asset.Length
	ctx.asset = asset
	//Save the length of the longest asset
	//Save the asset with the longest length
	// if ctx.length < asset.Length {
	// 	ctx.length = asset.Length
	// 	ctx.asset = asset
	// 	asset.AdjK = 0
	// } else if ctx.length > asset.Length {
	// 	//Longest asset is 134 length
	// 	//Our asset is 73
	// 	//The index when our asset 'start'
	// 	//is 134-73 = 61
	// 	asset.AdjK = ctx.length - asset.Length
	// 	fmt.Printf("Adjusted length is used: %d\n", asset.AdjK)
	// } else {
	// 	asset.AdjK = 0
	// }

}

func (ctx *Context) RemoveAsset(asset *D.Asset) {
	delete(ctx.assetMap, asset.Name)
	delete(ctx.assetIndicatorMap, asset.Name)
	delete(ctx.assetIndicatorMap, asset.Name)
}
func (ctx *Context) SetDataPath(path string) {
	ctx.dataPath = path
}

func (ctx *Context) GetDataPath() string {
	return ctx.dataPath
}
