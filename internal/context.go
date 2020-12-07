package internal

import (
	"time"

	log "github.com/sirupsen/logrus"
)

//Context for this backtester
type Context struct {
	Updateable   []Updateable
	Strategy     []Strategy
	Asset        []Asset
	AssetMap     map[string]*Asset
	StartDate    time.Time
	EndDate      time.Time
	Portfolio    Portfolio
	Broker       Broker
	K            int
	datePointer  time.Time
	eventChannel chan Event
	orderChannel chan Event
}

//NewContext creates a new context
func NewContext() *Context {
	eventChannelc := make(chan Event, 1)
	orderChannel := make(chan Event, 1)
	portfolio := Portfolio{}
	ctx := &Context{
		Asset:        make([]Asset, 1),
		AssetMap:     make(map[string]*Asset),
		eventChannel: eventChannelc,
		orderChannel: orderChannel,
		Portfolio:    portfolio,
		Broker:       Broker{portfolio: portfolio, channel: orderChannel},
	}
	log.WithFields(log.Fields{
		"Asset":        ctx.Asset,
		"AssetMap":     ctx.AssetMap,
		"eventChannel": ctx.eventChannel,
		"orderChannel": ctx.orderChannel,
		"Portfolio":    ctx.Portfolio,
		"Broker":       ctx.Broker,
	}).Debug("Created context")
	return ctx
}

func (ctx *Context) EventChannel() chan Event {
	return ctx.eventChannel
}

func (ctx *Context) OrderChannel() chan Event {
	return ctx.orderChannel
}

//Time returns the time
func (m *Context) Time() time.Time {
	return m.datePointer
}

//IncOneDay is used to step time forward
func (m *Context) IncOneDay() {
	old := m.datePointer
	m.datePointer = m.datePointer.AddDate(0, 0, 1)
	log.WithFields(log.Fields{
		"old": old,
		"new": m.datePointer,
	}).Debug("New day")
	m.shift()
}

func (m *Context) shift() {
	for i := range m.Asset {
		numberOfShifts, _ := m.Asset[i].Shift(m.datePointer)
		m.K += numberOfShifts
	}
	for _, element := range m.AssetMap {
		element.Shift(m.datePointer)
	}
}

//AddEndDate is used to set the strategy that will be run
func (m *Context) AddEndDate(endTime time.Time) {
	m.EndDate = endTime
}

//AddStartDate is used to set the strategy that will be run
func (m *Context) AddStartDate(startTime time.Time) {
	m.StartDate = startTime
	m.datePointer = startTime
}

//AddStrategy is used to set the strategy that will be run
func (m *Context) AddStrategy(strategy *Strategy) {
	m.Strategy = append(m.Strategy, *strategy)
}

//AddAsset is used to add assets that the strategy will use
func (m *Context) AddAsset(asset *Asset) {
	m.Asset = append(m.Asset, *asset)
	m.AssetMap[asset.Name] = asset
	log.WithFields(log.Fields{
		"Asset": asset,
	}).Debug("Adding asset to context")
}

//AddNamedAsset is used to add an asset and a reference
func (m *Context) AddNamedAsset(asset *Asset, name string) {
	m.Asset = append(m.Asset, *asset)
	m.AssetMap[name] = asset
}
