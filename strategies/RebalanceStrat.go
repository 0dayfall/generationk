package strategies

import (
	"log"
	"time"

	K "github.com/0dayfall/generationk"
	I "github.com/0dayfall/generationk/indicators"
)

//Strategy strategy
type RebalanceStrat struct {
	ROC100 []float64
}

//Setup is used to declare what indicators will be used
func (rmi *RebalanceStrat) Once(ctx *K.Context, ohlc *K.OHLC) error {

	//The Simple Moving Average length 50 periods, the ones from 0 to 50 will be registred in the array as well
	rmi.ROC100 = I.ROC100(ohlc.Close, 66)

	//If the init period is set PerBar will not be called until the InitPeriod is reached
	ctx.SetInitPeriod(66)

	return nil
}

func (rmi *RebalanceStrat) GetInterval() string { return "Q" }

var buytime time.Time

func (rmi *RebalanceStrat) Rebalance(k int, date time.Time, callback K.Callback) error {
	_, _, day := date.Date()

	if day == 28 {
		err := callback.SendOrder(K.BuyOrder, K.MarketOrder, 100)

		if err != nil {
			log.Fatal(err)

			return err
		}

		buytime = date
	}

	if callback.Owning() {
		if date.Sub(buytime).Hours()/24 > 66 {

			err := callback.SendOrder(K.SellOrder, K.MarketOrder, 100)

			if err != nil {
				log.Fatal(err)

				return err
			}
		}
	}

	return nil
}

//Tick get called when there is new data coming in
func (rmi *RebalanceStrat) PerBar(k int, callback K.Callback) error { return nil }

//Update is called before perBar
func (rmi *RebalanceStrat) Update(k *int) error { return nil }

//OrderEvent gets called on order events
func (rmi *RebalanceStrat) OrderEvent(orderEvent K.Event) {}
