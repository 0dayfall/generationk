⚠️ The package can be used for backtesting, but the API is not stable yet. Once stability is reached, version `1.0.0` will be tagged.
It is recommended your package manager locks to a commit id instead of the master branch directly. ⚠️

# Features

1. Read CVS file with historic data and back test

# ToDo

1. To set entry and exit conditions independently to be able to combine many different ones
2. To be able to try which parameters are the best ones for an indicator
3. 

# generationk
The inspiration for this project took place after using a few other backtesting frameworks in Python. 

I was tired of waiting for results and concluded that I want the fast feeling of a compiled language and I wanted all processor cores to be used.

I looked at a few different ones in Golang but they where either very complex for simple tasks or did not really appeal to the context.

## Design choices
The very first version was based on channels and real time but it very rarely stock data is real time in backtesting, even a minute is considered very granular and 5 minutes or 10 minutes are much more common. In that case reading a 2-3 data points will advance time 20-30 minutes just for placing an order which is not realistic in the
end so I decided to use callbacks instead.

## The Crossing MA example looks like this
```golang
type MACrossStrategy struct {
	ma50       *indicators.SimpleMovingAverage
	close      *indicators.TimeSeries
	initPeriod int
}

func (m *MACrossStrategy) Once(ctx *Context) error {
	m.close, e = NewTimeSeries(OHLC.Close, 5)
	m.ma50, e = NewSimpleMovingAverage(OHLC.Close, 50)
	if e != nil {
		return e
	}
	ctx.AddUpdatable(m.close, m.ma50)
}

func (m *MACrossStrategy) PerBar(callback Callback) {
	if m.ma50.Current() < m.close.Current() {
		if !callback.IsOwning("ABB") {
			callback.OrderSend("ABB", OrderDirection(Long), OrderType(Buy), 0, 1000)
		}
	}
}
```

There is also coded needed to create the strategy and run the backtest.
