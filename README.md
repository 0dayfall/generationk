⚠️ The package can be used for backtesting, but the API is not stable yet. Once stability is reached, version `1.0.0` will be tagged. It is recommended your package manager locks to a commit id instead of the master branch directly. ⚠️

# Features

1. Read CVS file with historic data and back test

# ToDo

1. To set entry and exit conditions independently to be able to combine many different ones
2. To be able to try which parameters are the best ones for an indicator
3. Many indicators can be implemented
4. Many candles can be implemented too
5. It would be interesting to output some common statistics about back test except profit etc.

# generationk
The inspiration for this project took place after using a few other backtesting frameworks in Python. 

I was tired of waiting for results and concluded that I want the fast feeling of a compiled language and I wanted all processor cores to be used.

I looked at a few different ones in Golang but they where either very complex for simple tasks or did not really appeal to the context.

## Design choices
I have created 3 versions of this framework. 
1. The very first version was based on channels and real time but it very rarely stock data is real time in backtesting, even a minute is considered very granular and 5 minutes or 10 minutes are much more common. In that case reading a 2-3 data points will advance time 20-30 minutes just for placing an order which is not realistic in the end.

2. The second version was based on callbacks using interfaces which would have the advantage of feeding one
data point after the other into the backtester from the simulated data source and this way obtain a more
realistic feeling. A main argument is that it is not possible to peek at future data by mistake. In the end
I think a backtest that peeks into future data stands out in the results as amazing and is actually easy to spot. I also felt that everything gets very slow and the reason is not strong enough since backtests are not
realistic anyway. My main priority is to check statistics to see if there is an edge and get a quick overview and then to see how I could manually trade a system.

3. In the third version which is the current what is provided is basically a for loop which loops over data. I have also tried to remove all kinds of encapsulation and special 'types'. I try to as much as possible only work with float64. Because it gives a degree of freedom and any functions can be written to work on numbers. So I want to as much as possible have access to the raw data and numbers instead of a 'LinearNumberSeries' struct and so on...

## The Crossing MA example looks like this
```golang
type MACrossStrategy struct {
	ma50  []float64
	close []float64
}

//Once is used to declare what indicators will be used; it's run once
func (ma *MACrossStrategy) Once(ctx *K.Context, ohlc *K.OHLC) error {

	//The closing prices
	ma.close = ohlc.Close

	//The Simple Moving Average length 50 periods, the ones from 0 to 50 
	//will be registred in the array as well
	ma.ma50 = indicators.SimpleMovingAverage(ohlc.Close, 50)

	//If the init period is set, the PerBar function below will not be called 
	//until the InitPeriod is reached
	ctx.SetInitPeriod(50)

	//No errors to my knowledge
	return nil
}

//PerBar gets called when there is new data coming in
func (ma *MACrossStrategy) PerBar(k int, callback K.Callback) error {

	//Closing price breaks the MA50
	if ma.close[k] > ma.ma50[k] {
		//Are we owning this stock since before?
		if !callback.Owning() {
			//No. Then we can buy 100 stocks with a marketorder
			callback.SendOrder(K.BuyOrder, K.MarketOrder, 100)

```
## Strategies
I want the strategies directory to contain working trading strategies that are actually used in the market. There are many working ones and the hard part is to trade the system, actually not find a system that is profitable. If you find something interesting, please share it. Sharing is caring 🤗

# Future implementation
* To use genetic algorithms to find the best trading system for a stocks / stocks
* To plug in machine learning