⚠️ This is an early version. The package can be used for backtesting, but the API is not stable yet. This is a work in progress. Once stability is reached, version `1.0.0` will be tagged. Please check back then. ⚠️

# Features

1. Read CVS file with historic data and back test in parallell. Testing ~5000 daily data per 100 CSV files can take about 1 sec. to complete.

# ToDo

1. To set entry and exit conditions independently to be able to combine many different ones, ex. entry RSI < 20; exit MA2 > MA5
2. To be able to try which parameters are the best ones for an indicator (simplified branch)
3. Use a genetic algo to scan the search space for most profitable parameters
4. Optimize for other functions than profit, drawdown, volatility
5. Make a WS API, deploy as lambda functions
6. Once a bit stable create a test parc 

# generationk
The inspiration for this project took place after using a few other backtesting frameworks in Python. 

I was tired of waiting for results and concluded that I want the fast feeling of a compiled language and I wanted all processor cores to be used.

I looked at a few different ones in Golang but they where either very complex for simple tasks or did not really appeal to the context.

## Design choices
The very first version was based on channels and real time but it very rarely stock data is real time in backtesting, even a minute is considered very granular and 5 minutes or 10 minutes are much more common. In that case reading a 2-3 data points will advance time 20-30 minutes just for placing an order which is not realistic in the
end so I decided to use callbacks instead.

## The Crossing MA example looks like this
```golang
//Strategy strategy
type MACrossStrategy struct {
	ma50  indicators.SimpleMovingAverage
	close indicators.TimeSeries
}

//Setup is used to declare what indicators will be used
func (ma *MACrossStrategy) Once(ctx *K.Context) error {
	//Want access to the latest 5 closing prices
	ma.close = indicators.NewTimeSeries(indicators.Close, 5)
	//MA50
	ma.ma50 = indicators.NewSimpleMovingAverage(indicators.Close, 50)

	//Add indicators to context
	ctx.AddIndicator(&ma.close)
	ctx.AddIndicator(&ma.ma50)

	//The data needed to calculate MA
	ctx.SetInitPeriod(50)

	return nil
}

//Tick get called when there is new data coming in
func (ma *MACrossStrategy) PerBar(callback K.Callback) {

	if ma.close.Current() > ma.ma50.Current() {
		if !callback.IsOwning("ABB") {
			err := callback.OrderSend("ABB", K.BuyOrder, K.MarketOrder, 0, 100)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if ma.close.Current() < ma.ma50.Current() {
		if callback.IsOwning("ABB") {
			err := callback.OrderSend("ABB", K.SellOrder, K.MarketOrder, 0, 100)
```

There is also coded needed to create the strategy and run the backtest.
