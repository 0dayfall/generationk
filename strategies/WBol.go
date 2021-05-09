package generationk

import (
	K "github.com/0dayfall/generationk"
	D "github.com/0dayfall/generationk/data"
	I "github.com/0dayfall/generationk/indicators"
)

//Strategy strategy
type WBol struct {
	up    map[*D.Asset][]float64
	dn    map[*D.Asset][]float64
	close map[*D.Asset][]float64
}

func (wbol *WBol) GetParams() []*K.Params { return nil }

//Setup is used to declare what indicators will be used
func (wbol *WBol) Once(ctx *K.Context, assets []*D.Asset) error {
	wbol.up = make(map[*D.Asset][]float64, len(assets))
	wbol.dn = make(map[*D.Asset][]float64, len(assets))
	wbol.close = make(map[*D.Asset][]float64, len(assets))

	for _, asset := range assets {
		ohlc := asset.Ohlc

		//The rate of change for the last 66 days
		wbol.close[asset] = ohlc.Close

		//The down, up
		wbol.up[asset], _, wbol.up[asset] = I.Bollinger(ohlc.Close, 2, 60)
	}

	//If the init period is set PerBar will not be called until the InitPeriod is reached
	ctx.SetInitPeriod(200)

	return nil
}

//Tick get called when there is new data coming in
func (wbol *WBol) PerBar(k int, callback K.Callback) error {
	//It happens per bar

	/*for _, asset := range  {

	}*/

	return nil
}

func (wbol *WBol) End(k int, callback K.Callback) error {
	//When it finish

	return nil
}

//Update is called before perBar
func (wbol *WBol) Update(k *int) error { return nil }

//OrderEvent gets called on order events
func (wbol *WBol) OrderEvent(orderEvent K.Event) {}
