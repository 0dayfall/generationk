# generationk
The inspiration for this project took place after using a few other backtesting frameworks in Python. I was tired of waiting for results and concluded that I want the fast feeling of a compiled language and I also want to make use the multiple processor cores that often is available but rarely used.

## Design choices
Going for event driven in the end came more from a point of being able to split the program up to run on multiple computers for performance reasons rather than a needs for realism (but implictly contributes to realism). Another choice was working with float arrays as arguments to all indicators to keep it as simple as possible.

## The Crossing MA example looks like this
```golang
type MACrossStrategy struct {
	ma50       *indicators.SimpleMovingAverage
	close      *indicators.TimeSeries
	initPeriod int
}

func (m *MACrossStrategy) Setup(ctx *Context) error {
	m.close, e = NewTimeSeries(ctx.AssetMap["AAPL"])
	m.ma50, e = NewSimpleMovingAverage(ctx.AssetMap["AAPL"], Close, 5)
	if e != nil {
		return e
	}
	ctx.AddUpdatable(m.close, m.ma50)
}

func (m *MACrossStrategy) Tick(ctx *Context) {
	if m.ma50.ValueAtIndex(0) > m.close.ValueAtIndex(0) {
		if !ctx.Position(ctx.AssetMap["ABB"]) {
			MakeOrder(ctx, OrderType(Buy), ctx.AssetMap["ABB"], ctx.Time(), 1000)
		}
	}
}
```

There is also coded needed to create the strategy and run the backtest.
