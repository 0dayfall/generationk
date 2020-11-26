package internal

import (
	"time"
)

//Context for this backtester
type Context struct {
	Strategy    []Strategy
	Asset       []Asset
	AssetMap    map[string]*Asset
	StartDate   time.Time
	EndDate     time.Time
	datePointer time.Time
}

//NewContext creates a new context
func NewContext() *Context {
	return &Context{
		AssetMap: make(map[string]*Asset),
	}
}

//Time returns the time
func (m *Context) Time() time.Time {
	return m.datePointer
}

//IncOneDay is used to step time forward
func (m *Context) IncOneDay() {
	m.datePointer.AddDate(0, 0, 1)
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
}

//AddNamedAsset is used to add an asset and a reference
func (m *Context) AddNamedAsset(asset *Asset, name string) {
	m.Asset = append(m.Asset, *asset)
	m.AssetMap[name] = asset
}
