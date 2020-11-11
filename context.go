package main

//Context for this backtester
type Context struct {
	strategy Strategy
	asset    Asset
}

func (m *Context) setStrategy(strategy *Strategy) {
	m.strategy = *strategy
}

func (m *Context) setAsset(asset *Asset) {
	m.asset = *asset
}
