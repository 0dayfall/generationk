package generationk

import (
	"fmt"
	"log"

	K "github.com/0dayfall/generationk"
	D "github.com/0dayfall/generationk/data"
	I "github.com/0dayfall/generationk/indicators"
)

//Strategy strategy
type RMIHigh struct {
	fastRMI   []float64
	slowRMI   []float64
	ema10     []float64
	close     []float64
	buyParam  K.Params
	sellParam K.Params
}

func (donch *RMIHigh) GetParams() []*K.Params { return nil }

//Setup is used to declare what indicators will be used
func (rmi *RMIHigh) Once(ctx *K.Context, asset []*D.Asset) error {
	ohlc := asset[0].Ohlc

	//The Simple Moving Average length 50 periods, the ones from 0 to 50 will be registred in the array as well
	rmi.fastRMI, rmi.slowRMI = I.RMI(ohlc.Close, 15, 2, 250, 20)
	rmi.ema10 = I.EMA(ohlc.Close, 10)
	rmi.close = ohlc.Close

	//If the init period is set PerBar will not be called until the InitPeriod is reached
	ctx.SetInitPeriod(50)

	return nil
}

//Tick get called when there is new data coming in
func (rmi *RMIHigh) PerBar(k int, callback K.Callback) error {

	if rmi.fastRMI[k] > 27 {
		if !callback.Owning() {
			fmt.Println("Buying")
			_, err := callback.SendOrder(K.BuyOrder, K.MarketOrder, 100)

			if err != nil {
				log.Fatal(err)

				return err
			}
		}
	}

	if rmi.close[k] < rmi.ema10[k] {
		if callback.Owning() {
			fmt.Println("Selling")
			_, err := callback.SendOrder(K.SellOrder, K.MarketOrder, 100)

			if err != nil {
				log.Fatal(err)

				return err
			}
		}
	}

	return nil
}

func (rmi *RMIHigh) End(k int, callback K.Callback) error { return nil }

//Update is called before perBar
func (rmi *RMIHigh) Update(k *int) error { return nil }

//OrderEvent gets called on order events
func (rmi *RMIHigh) OrderEvent(orderEvent K.Event) {}
