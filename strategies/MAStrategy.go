package generationk

import (
	"log"

	K "github.com/0dayfall/generationk"
	D "github.com/0dayfall/generationk/data"
	I "github.com/0dayfall/generationk/indicators"
)

//Strategy strategy
type MACrossStrategy struct {
	ma50  []float64
	close []float64
}

func (ma *MACrossStrategy) GetParams() []*K.Params {
	return nil
}

//Setup is used to declare what indicators will be used
func (ma *MACrossStrategy) Once(ctx *K.Context, asset []*D.Asset) error {

	ohlc := asset[0].Ohlc
	//The closing prices
	ma.close = ohlc.Close

	//The Simple Moving Average length 50 periods, the ones from 0 to 50 will be registred in the array as well
	ma.ma50 = I.SimpleMovingAverage(ohlc.Close, 50)

	//If the init period is set PerBar will not be called until the InitPeriod is reached
	ctx.SetInitPeriod(50)

	return nil
}

//Tick get called when there is new data coming in
func (ma *MACrossStrategy) PerBar(k int, callback K.Callback) error {

	/*fmt.Printf("k: %d ", k)
	fmt.Printf("len(close): %d ", len(ma.close))
	fmt.Printf("len(ma5) %d\n", len(ma.ma50))*/

	if ma.close[k] > ma.ma50[k] {
		if !callback.Owning() {
			_, err := callback.SendOrder(K.BuyOrder, K.MarketOrder, 100)

			if err != nil {
				log.Fatal(err)

				return err
			}
		}
	}

	if ma.close[k] < ma.ma50[k] {
		if callback.Owning() {
			_, err := callback.SendOrder(K.SellOrder, K.MarketOrder, 100)

			if err != nil {
				log.Fatal(err)

				return err
			}
		}
	}

	return nil
}

func (ma *MACrossStrategy) End(k int, callback K.Callback) error { return nil }

//Update is called before perBar
func (ma *MACrossStrategy) Update(k *int) error { return nil }

//OrderEvent gets called on order events
func (ma *MACrossStrategy) OrderEvent(orderEvent K.Event) {}
