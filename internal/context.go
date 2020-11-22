package internal

import (
	"fmt"
	"time"
)

//Context for this backtester
type Context struct {
	Strategy    []Strategy
	Asset       []Asset
	StartDate   time.Time
	EndDate     time.Time
	datePointer time.Time
}

//NewContext creates a new context
func (m *Context) NewContext() *Context {
	return &Context{}
}

//IncTime is used to step time forward
func (m *Context) IncTime() {
	m.datePointer.AddDate(0, 0, 1)
	for i := range m.Asset {
		//asset.Ohlc = asset.Ohlc.shift()
		//_, asset.Ohlc = asset.Ohlc[0], asset.Ohlc[1:]
		m.Asset[i].Ohlc = m.Asset[i].Ohlc[1:]
		fmt.Printf("New value: %f\n", m.Asset[i].Ohlc[1:][0].Close)
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
