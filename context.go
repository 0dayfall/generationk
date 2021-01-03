package generationk

import (
	"time"

	indicators "github.com/0dayfall/generationk/indicators"
	log "github.com/sirupsen/logrus"
)

//Context holds holds the strategy, assets, indicators per asset, start date, end date, the portfolio
//the current date and the unstable period
type Context struct {
	strategy          []Strategy
	assets            []Asset
	assetMap          map[string]*Asset
	assetIndicatorMap map[string][]indicators.Indicator
	startDate         time.Time
	endDate           time.Time
	portfolio         *Portfolio
	broker            Broker
	K                 int
	datePointer       time.Time
	/*eventChannel      chan Event
	orderChannel      chan Event*/
	initPeriod int
}

//NewContext creates a new context
func NewContext() *Context {
	//eventChannelc := make(chan Event, 1)
	//orderChannel := make(chan Event, 1)
	ctx := &Context{
		assets:            []Asset{},
		assetMap:          make(map[string]*Asset),
		assetIndicatorMap: make(map[string][]indicators.Indicator),
		//eventChannel:      eventChannelc,
		//orderChannel:      orderChannel,
		broker: Broker{},
	}

	log.WithFields(log.Fields{
		"Asset":    ctx.assets,
		"AssetMap": ctx.assetMap,
		/*		"eventChannel": ctx.eventChannel,
				"orderChannel": ctx.orderChannel,*/
		"Portfolio": &ctx.portfolio,
		"Broker":    ctx.broker,
	}).Info("Created context")

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
	for name, asset := range ctx.assetMap {
		ctx.assetIndicatorMap[name] = append(ctx.assetIndicatorMap[name], indicator)

		log.WithFields(log.Fields{
			"ctx.AssetIndicatorMap[k]":    indicator.GetName(),
			"ctx.AssetMap":                asset.name,
			"ctx.AssetIndicatorMap[name]": ctx.assetIndicatorMap[name],
		}).Info("Adding indicator to asset")
		/*if indicator.GetDataType() == indicators.Close {
			indicator.Update(asset.CloseArray())
		}*/
	}
}

//AddEndDate is used to set the end date
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

	log.WithFields(log.Fields{
		"period": ctx.initPeriod,
	}).Info("Setting the init period")
}

//GetInitPeriod returns the period
func (ctx *Context) GetInitPeriod() int {
	log.WithFields(log.Fields{
		"period": ctx.initPeriod,
	}).Debug("Returning the init period")

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
	ctx.assets = append(ctx.assets, *asset)
	ctx.assetMap[asset.name] = asset
	ctx.assetIndicatorMap[asset.name] = nil

	log.WithFields(log.Fields{
		"Asset": asset.name,
	}).Debug("Asset added to context")
}

//updateIndicators is used by generationK to update the data
func (ctx *Context) updateIndicators(assetName string) {
	log.Debug("ctx.AssetIndicatorMap[assetName]: ", len(ctx.GetAssetIndicatorByName(assetName)))

	for _, indicator := range ctx.GetAssetIndicatorByName(assetName) {

		//Copy period amount of data to update indicator with
		dataSlice := ctx.GetAssetByName(assetName).GetData(
			OhlcConst(indicator.GetDataType()),
			indicator.GetPeriod(),
		)
		indicator.Update(dataSlice)
	}
}
