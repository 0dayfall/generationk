package generationk

import (
	"fmt"
	"log"
	"strings"

	K "github.com/0dayfall/generationk"
	D "github.com/0dayfall/generationk/data"
	I "github.com/0dayfall/generationk/indicators"
)

//Strategy strategy
type DonchStrategy struct {
	up     map[*D.Asset][]float64
	dn     map[*D.Asset][]float64
	close  map[*D.Asset][]float64
	ma20   map[*D.Asset][]float64
	ma50   map[*D.Asset][]float64
	ma100  map[*D.Asset][]float64
	ma200  map[*D.Asset][]float64
	roc100 map[*D.Asset][]float64
	bull   bool
	trend  float64
}

func (donch *DonchStrategy) GetParams() []*K.Params { return nil }

//Setup is used to declare what indicators will be used
func (donch *DonchStrategy) Once(ctx *K.Context, assets []*D.Asset) error {
	donch.up = make(map[*D.Asset][]float64, len(assets))
	donch.dn = make(map[*D.Asset][]float64, len(assets))
	donch.close = make(map[*D.Asset][]float64, len(assets))
	donch.ma20 = make(map[*D.Asset][]float64, len(assets))
	donch.ma50 = make(map[*D.Asset][]float64, len(assets))
	donch.ma100 = make(map[*D.Asset][]float64, len(assets))
	donch.ma200 = make(map[*D.Asset][]float64, len(assets))
	donch.roc100 = make(map[*D.Asset][]float64, len(assets))

	for _, asset := range assets {
		ohlc := asset.Ohlc

		//The rate of change for the last 66 days
		donch.close[asset] = ohlc.Close
		//The Donchian channel length x periods
		donch.up[asset], _, donch.dn[asset] = I.Donchian(ohlc.High, ohlc.Low, 100)
		donch.ma20[asset] = I.MA(ohlc.Close, 20)
		donch.ma50[asset] = I.MA(ohlc.Close, 50)
		donch.ma100[asset] = I.MA(ohlc.Close, 100)
		donch.ma200[asset] = I.MA(ohlc.Close, 200)
		donch.roc100[asset] = I.ROC100(ohlc.Close, 100)
	}

	//If the init period is set PerBar will not be called until the InitPeriod is reached
	ctx.SetInitPeriod(200)

	return nil
}

//Tick get called when there is new data coming in
func (donch *DonchStrategy) PerBar(k int, callback K.Callback) error {

	//assetCount := 0
	for asset := range donch.ma200 {
		if strings.ToLower(asset.Name) == "omx" {
			if donch.close[asset][k-asset.AdjK] > donch.ma100[asset][k-asset.AdjK-1] {
				donch.bull = true
			} else {
				donch.bull = false
			}
		}
		/*assetCount++
		if donch.close[asset][k-asset.AdjK] > donch.ma200[asset][k-asset.AdjK] {
			donch.trend++
		}*/

	}

	/*donch.trend = donch.trend / float64(assetCount)
	fmt.Println(donch.trend)*/

	/*if donch.trend > 0.7 {
		donch.bull = true
	} else {
		donch.bull = false
	}*/

	for asset := range donch.up {
		//fmt.Printf("Close %f donch %f\n", donch.close[asset][k-asset.AdjK], donch.up[asset][k-asset.AdjK])
		close := donch.close[asset][k-asset.AdjK]

		if close > donch.up[asset][k-asset.AdjK-1] &&
			donch.roc100[asset][k-asset.AdjK] > 30 &&
			donch.bull {

			owning, err := callback.IsOwning(asset.Name)
			if err != nil {
				return err
			}

			if !owning {

				cost, err := callback.SendOrderFor(asset.Name, K.BuyOrder, K.MarketOrder, callback.PositionSize(0.1, close))
				if err != nil {
					//log.Fatal(err)

					return nil
				}

				callback.Record(asset.Name, K.RecordStruct{
					Time: callback.Date(),
					Variable: map[string]string{
						"Donh": fmt.Sprintf("%f", donch.up[asset][k-1]),
						"Cost": fmt.Sprintf("%f", cost),
					}})
			}
		}

		if donch.close[asset][k-asset.AdjK] < donch.ma20[asset][k-asset.AdjK] &&
			donch.close[asset][k-asset.AdjK-1] < donch.ma20[asset][k-asset.AdjK-1] {

			owning, err := callback.IsOwning(asset.Name)
			if err != nil {
				return err
			}

			if owning {

				cash, err := callback.Sell(asset.Name)
				if err != nil {
					log.Fatal(err)

					return err
				}

				callback.Record(asset.Name, K.RecordStruct{
					Time: callback.Date(),
					Variable: map[string]string{
						"Donh": fmt.Sprintf("%f", donch.dn[asset][k-1]),
						"Cash": fmt.Sprintf("%f", cash),
					}})
			}
		}
	}
	return nil
}

func (donch *DonchStrategy) End(k int, callback K.Callback) error { return nil }

//Update is called before perBar
func (donch *DonchStrategy) Update(k *int) error { return nil }

//OrderEvent gets called on order events
func (donch *DonchStrategy) OrderEvent(orderEvent K.Event) {}
