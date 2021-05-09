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
const HoldingDays = 66

var balance []float64

//Strategy strategy
type RebalanceStrat struct {
	ROC66   map[*D.Asset][]float64
	ROC22   map[*D.Asset][]float64
	time    map[*D.Asset][]time.Time
	buyTime map[string]time.Time
}

func (reb *RebalanceStrat) GetParams() []*K.Params {
	return nil
}

//Setup is used to declare what indicators will be used
func (reb *RebalanceStrat) Once(ctx *K.Context, assets []*D.Asset) error {
	reb.time = make(map[*D.Asset][]time.Time, len(assets))
	reb.ROC66 = make(map[*D.Asset][]float64, len(assets))
	reb.ROC22 = make(map[*D.Asset][]float64, len(assets))

	reb.buyTime = make(map[string]time.Time)

	for _, asset := range assets {
		ohlc := asset.Ohlc

		//The rate of change for the last 66 days
		reb.time[asset] = ohlc.Time
		reb.ROC66[asset] = I.ROC100(ohlc.Close, 120)
		reb.ROC22[asset] = I.ROC100(ohlc.Close, 22)
	}

	//If the init period is set PerBar will not be called until the InitPeriod is reached
	ctx.SetInitPeriod(HoldingDays)

	return nil
}

func (reb *RebalanceStrat) GetInterval() string { return "Q" }

type roc struct {
	name  string
	asset *D.Asset
	value float64
}

func (reb *RebalanceStrat) Rebalance(k int, callback K.Callback) error {
	_, _, day := callback.Date().Date()
	if day == 28 {
		keys66 := make([]roc, 0, len(reb.ROC66))

		for asset, close := range reb.ROC66 {
			if k-asset.AdjK > 0 {
				if k-asset.AdjK < len(close)-1 {
					keys66 = append(keys66, roc{asset.Name, asset, close[k-asset.AdjK]})
				}
			}
		}

		sort.Slice(keys66, func(i, j int) bool {
			return keys66[i].value > keys66[j].value
		})

		for i := 0; i < Holdings; i++ {

			//if keys66[i].value > 20 && reb.ROC22[keys66[i].asset][k-keys66[i].asset.AdjK] > 5.0 {
			if keys66[i].value > 20 {

				cost, err := callback.SendOrderFor(keys66[i].name, K.BuyOrder, K.MarketOrder, callback.PositionSize(0.02, keys66[i].value))
				if err != nil {
					log.Fatal(err)

					return err
				}

				callback.Record(keys66[i].name, K.RecordStruct{
					Time: callback.Date(),
					Variable: map[string]string{
						"ROC Value": fmt.Sprintf("%f", keys66[i].value),
						"Cost":      fmt.Sprintf("%f", cost),
					}})

				balance = append(balance, cost)
				reb.buyTime[keys66[i].name] = callback.Date()
			}
		}
	}

	return nil
}

//Tick get called when there is new data coming in
func (reb *RebalanceStrat) PerBar(k int, callback K.Callback) error {
	date := callback.Date()

	for asset := range reb.ROC66 {

		owning, err := callback.IsOwning(asset.Name)
		if err != nil {
			return err
		}

		if owning {

			if timeDiff, ok := reb.buyTime[asset.Name]; ok {

				if date.Sub(timeDiff).Hours()/24 > HoldingDays {

					cash, err := callback.SendOrderFor(asset.Name, K.SellOrder, K.MarketOrder, callback.PositionSize(0.02, asset.Ohlc.Close[k]))
					if err != nil {
						return err
					}

					callback.Record(asset.Name, K.RecordStruct{
						Time: callback.Date(),
						Variable: map[string]string{
							"Cash:": fmt.Sprintf("%f", cash),
						}})

					balance = append(balance, cash)
					delete(reb.buyTime, asset.Name)

				}
			}
		}

	}

	return nil
}

func (reb *RebalanceStrat) End(k int, callback K.Callback) error {

	for asset := range reb.ROC66 {

		owning, err := callback.IsOwning(asset.Name)
		if err != nil {
			return err
		}

		if owning {

			cash, err := callback.SendOrderFor(asset.Name, K.SellOrder, K.MarketOrder, callback.PositionSize(0.02, asset.Ohlc.Close[k]))
			if err != nil {
				return err
			}

			balance = append(balance, cash)
		}
	}

	fmt.Println(balance)

	var sum float64
	for _, value := range balance {
		sum += value
	}
	fmt.Printf("\nSum: %f\n", sum)

	return nil
}

//Update is called before perBar
func (reb *RebalanceStrat) Update(k *int) error { return nil }

//OrderEvent gets called on order events
func (reb *RebalanceStrat) OrderEvent(orderEvent K.Event) {}
