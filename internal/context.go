package internal

import ("time")

//Context for this backtester
type Context struct {
	Strategy []Strategy
	Asset    []Asset
	StartDate time.Time
	EndDate time.Time
}

//AddEndDate is used to set the strategy that will be run
func (m *Context) AddEndDate(endTime time.Time) {
	m.EndDate = endTime
}

//AddStartDate is used to set the strategy that will be run
func (m *Context) AddStartDate(startTime time.Time) {
	m.StartDate = startTime
}

//AddStrategy is used to set the strategy that will be run
func (m *Context) AddStrategy(strategy *Strategy) {
	m.Strategy = append(m.Strategy, *strategy)
}

//AddAsset is used to add assets that the strategy will use
func (m *Context) AddAsset(asset *Asset) {
	m.Asset = append(m.Asset, *asset)
}
