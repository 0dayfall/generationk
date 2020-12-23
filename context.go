package generationk

import (
	"time"

	indicators "github.com/0dayfall/generationk/indicators"
	log "github.com/sirupsen/logrus"
)

//Context for this backtester
type Context struct {
	//Updateable        []Updateable
	Strategy          []Strategy
	Asset             []Asset
	AssetMap          map[string]*Asset
	AssetIndicatorMap map[string][]*indicators.Indicator
	StartDate         time.Time
	EndDate           time.Time
	Portfolio         *Portfolio
	Broker            Broker
	K                 int
	datePointer       time.Time
	eventChannel      chan Event
	orderChannel      chan Event
	period            int
}

//NewContext creates a new context
func newContext() *Context {
	//eventChannelc := make(chan Event, 1)
	//orderChannel := make(chan Event, 1)

	ctx := &Context{
		Asset:             make([]Asset, 1),
		AssetMap:          make(map[string]*Asset),
		AssetIndicatorMap: make(map[string][]*indicators.Indicator),
		//eventChannel:      eventChannelc,
		//orderChannel:      orderChannel,
		Broker: Broker{},
	}
	log.WithFields(log.Fields{
		"Asset":        ctx.Asset,
		"AssetMap":     ctx.AssetMap,
		"eventChannel": ctx.eventChannel,
		"orderChannel": ctx.orderChannel,
		"Portfolio":    &ctx.Portfolio,
		"Broker":       ctx.Broker,
	}).Info("Created context")
	return ctx
}

//EventChannel return the channel for events
func (ctx *Context) EventChannel() chan Event {
	return ctx.eventChannel
}

//OrderChannel returns the channel for orders
func (ctx *Context) OrderChannel() chan Event {
	return ctx.orderChannel
}

//Time returns the time
func (ctx *Context) Time() time.Time {
	return ctx.datePointer
}

//IncOneDay is used to step time forward
/*func (ctx *Context) IncOneDay() {
	old := ctx.datePointer
	ctx.datePointer = ctx.datePointer.AddDate(0, 0, 1)
	log.WithFields(log.Fields{
		"old": old,
		"new": ctx.datePointer,
	}).Debug("New day")
	ctx.shift()
}*/

//AddIndicatorOnAsset will add an indicator on the asset
func (ctx *Context) AddIndicatorOnAsset(asset *Asset, indicator *indicators.Indicator) {
	ctx.AssetIndicatorMap[asset.Name] = append(ctx.AssetIndicatorMap[asset.Name], indicator)
}

//AddIndicator will add it to all assets
func (ctx *Context) AddIndicator(indicator indicators.Indicator) {
	for name, asset := range ctx.AssetMap {
		ctx.AssetIndicatorMap[name] = append(ctx.AssetIndicatorMap[name], &indicator)
		log.WithFields(log.Fields{
			"ctx.AssetIndicatorMap[k]": indicator.GetName(),
			"ctx.AssetMap":             asset.Name,
		}).Debug("Adding indicator to asset")
		if indicator.GetDataType() == indicators.Close {
			indicator.Update(ctx.AssetMap[name].CloseArray())
		}
	}
}

//AddUpdatable add an updatable interface
func (ctx *Context) AddUpdatable(indicators ...Updateable) {
	//ctx.Updateable = indicators
}

/*func (ctx *Context) shift() {
	for i := range ctx.Asset {
		numberOfShifts, _ := ctx.Asset[i].Shift(ctx.datePointer)
		ctx.K += numberOfShifts
	}
	for _, element := range ctx.AssetMap {
		element.Shift(ctx.datePointer)
	}
}*/

//AddEndDate is used to set the strategy that will be run
func (ctx *Context) AddEndDate(endTime time.Time) {
	ctx.EndDate = endTime
}

//AddStartDate is used to set the strategy that will be run
func (ctx *Context) AddStartDate(startTime time.Time) {
	ctx.StartDate = startTime
	ctx.datePointer = startTime
}

//AddStrategy is used to set the strategy that will be run
func (ctx *Context) AddStrategy(strategy *Strategy) {
	ctx.Strategy = append(ctx.Strategy, *strategy)
}

//SetInitPeriod is used to set the priod for which
func (ctx *Context) SetInitPeriod(period int) {
	ctx.period = period
	log.WithFields(log.Fields{
		"period": ctx.period,
	}).Debug("Setting the init period")
}

//GetInitPeriod returns the period
func (ctx Context) GetInitPeriod() int {
	log.WithFields(log.Fields{
		"period": ctx.period,
	}).Debug("Returning the init period")
	return ctx.period
}

//AddAsset is used to add assets that the strategy will use
func (ctx *Context) AddAsset(asset *Asset) {
	ctx.Asset = append(ctx.Asset, *asset)
	ctx.AssetMap[asset.Name] = asset
	ctx.AssetIndicatorMap[asset.Name] = nil

	log.WithFields(log.Fields{
		"Asset": asset.Name,
	}).Debug("Asset added to context")
}
