package internal

import ("time")

//Context for this backtester
type Context struct {
	Strategy []Strategy
	Asset    []Asset
	StartDate time.Time
	EndDate time.Time
	datePointer time.Time
}

//NewContext creates a new context
func (m *Context) NewContext() *Context {
	return &Context{Strategy: nil, Asset: nil, StartDate: time.Now(), EndDate: time.Now(), datePointer: time.Now()}
}

//IncTime is used to step time forward
func (m *Context) IncTime() {
	m.datePointer.AddDate(0,0,1)
	for _, asset := range m.Asset {
		//asset.Ohlc = asset.Ohlc.shift()
		_, asset.Ohlc = asset.Ohlc[0], asset.Ohlc[1:]
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
}
