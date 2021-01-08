package generationk

//Strategy is the class where the logic is placed to buy and sell assets
//the two methods that needs to be implemented are Setup and Tick.
//The Setup method is used to define if any indicators will be used
//and what period they need to be stable.
//The Tick method is called for every new data which arrives and is
//a possibility to make checks and send orders.
type Strategy interface {
	Once(ctx *Context, ohlc *OHLC) error
	Update(k *int) error
	PerBar(k int, callback Callback)
}

type MultiStrategy interface {
	Once(ctx *Context, ohlc OHLC) error
	PerBar(k *int, callback Callback)
}
