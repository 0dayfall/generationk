package internal

//Event type
type Event struct{}

//Tick event type
type Tick Event

//Signal event type
type Signal Event

//Order event type
type Order Event

//Fill event type
type Fill Event

//Data event type
type Data Event

//PutData starts a backtest with the information in context
func PutData(c *Context) {
	data := make(chan Data)

	for i, asset := range c.Assets {
		for j, ohlc := range asset[j] {
			data <- asset[i].ohlc[j]
		}
	}
}

//Run starts a backtest with the information in context
func GetData(c *Context) {
	data := make(chan Data)

	for i, asset := range c.Assets {
		for j, ohlc := range asset[j] {
			data <- asset[i].ohlc[j]
		}
	}
}
