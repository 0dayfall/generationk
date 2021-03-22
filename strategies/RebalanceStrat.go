package generationk

import (
	"fmt"
	"log"
	"sort"
	"time"

	K "github.com/0dayfall/generationk"
	D "github.com/0dayfall/generationk/data"
	I "github.com/0dayfall/generationk/indicators"
)

const Holdings = 5
const HoldingDays = 45

//Strategy strategy
type RebalanceStrat struct {
	ROC100  map[*D.Asset][]float64
	time    map[*D.Asset][]time.Time
	buyTime map[string]time.Time
}

func (rmi *RebalanceStrat) GetParams() []*K.Params {
	return nil
}

//Setup is used to declare what indicators will be used
func (rmi *RebalanceStrat) Once(ctx *K.Context, assets []*D.Asset) error {
	rmi.time = make(map[*D.Asset][]time.Time, len(assets))
	rmi.ROC100 = make(map[*D.Asset][]float64, len(assets))
	rmi.buyTime = make(map[string]time.Time)

	for _, asset := range assets {
		ohlc := asset.Ohlc

		//The rate of change for the last 66 days
		rmi.time[asset] = ohlc.Time
		rmi.ROC100[asset] = I.ROC100(ohlc.Close, 66)
	}
	//If the init period is set PerBar will not be called until the InitPeriod is reached
	ctx.SetInitPeriod(66)

	return nil
}

func (rmi *RebalanceStrat) GetInterval() string { return "Q" }

type roc struct {
	name  string
	value float64
}

func (rmi *RebalanceStrat) Rebalance(k int, callback K.Callback) error {
	_, _, day := callback.Date().Date()
	if day == 28 {
		fmt.Print("\n", callback.Date())
		keys := make([]roc, 0, len(rmi.ROC100))

		for asset, close := range rmi.ROC100 {
			if k-asset.AdjK > 0 {
				if k-asset.AdjK < len(close)-1 {
					keys = append(keys, roc{asset.Name, close[k-asset.AdjK]})
				}
			}
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i].value > keys[j].value
		})

		for i := 0; i < Holdings; i++ {
			fmt.Printf("\n%s, %f, ", keys[i].name, keys[i].value)
			if keys[i].value > 20 {
				err := callback.SendOrderFor(keys[i].name, K.BuyOrder, K.MarketOrder, 100)

				if err != nil {
					log.Fatal(err)

					return err
				}
				rmi.buyTime[keys[i].name] = callback.Date()
			}
		}
		fmt.Printf("\n")
	}

	return nil
}

//Tick get called when there is new data coming in
func (rmi *RebalanceStrat) PerBar(k int, callback K.Callback) error {
	date := callback.Date()

	for asset, _ := range rmi.ROC100 {
		owning, err := callback.IsOwning(asset.Name)
		if err != nil {
			return err
		}

		if owning {
			if timeDiff, ok := rmi.buyTime[asset.Name]; ok {
				if date.Sub(timeDiff).Hours()/24 > HoldingDays {
					fmt.Printf("Bought=> %s Now=> %s", timeDiff, date)
					fmt.Printf("Held %s > 66 days\n", asset.Name)
					err := callback.SendOrderFor(asset.Name, K.SellOrder, K.MarketOrder, 100)
					if err != nil {
						return err
					}
					delete(rmi.buyTime, asset.Name)
				}
			}
		}
	}
	return nil
}

//Update is called before perBar
func (rmi *RebalanceStrat) Update(k *int) error { return nil }

//OrderEvent gets called on order events
func (rmi *RebalanceStrat) OrderEvent(orderEvent K.Event) {}
