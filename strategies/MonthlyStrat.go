package generationk

import (
	"fmt"
	"log"
	"time"

	K "github.com/0dayfall/generationk"
	D "github.com/0dayfall/generationk/data"
	I "github.com/0dayfall/generationk/indicators"
)

//Strategy strategy
type MonthlyStrat struct {
	up     map[*D.Asset][]float64
	dn     map[*D.Asset][]float64
	close  map[*D.Asset][]float64
	roc60  map[*D.Asset][]float64
	roc120 map[*D.Asset][]float64
	bull   bool
	//trend  float64
	oldMonth     int64
	rebalanceDay int64
}

func (mrebstrat *MonthlyStrat) GetParams() []*K.Params { return nil }

//Setup is used to declare what indicators will be used
func (mrebstrat *MonthlyStrat) Once(ctx *K.Context, assets []*D.Asset) error {
	mrebstrat.close = make(map[*D.Asset][]float64, len(assets))
	mrebstrat.up = make(map[*D.Asset][]float64, len(assets))
	mrebstrat.dn = make(map[*D.Asset][]float64, len(assets))
	mrebstrat.roc60 = make(map[*D.Asset][]float64, len(assets))
	mrebstrat.roc120 = make(map[*D.Asset][]float64, len(assets))

	for _, asset := range assets {
		ohlc := asset.Ohlc

		//The rate of change for the last 66 days
		mrebstrat.close[asset] = ohlc.Close

		//The mrebstratian channel length x periods
		mrebstrat.up[asset], _, mrebstrat.dn[asset] = I.Donchian(ohlc.High, ohlc.Low, 100)
		mrebstrat.roc60[asset] = I.ROC100(ohlc.Close, 60)
		mrebstrat.roc120[asset] = I.ROC100(ohlc.Close, 120)
	}

	mrebstrat.oldMonth = 0
	mrebstrat.rebalanceDay = 28

	//If the init period is set PerBar will not be called until the InitPeriod is reached
	ctx.SetInitPeriod(200)

	return nil
}

// A weekend is Friday 10pm UTC to Sunday 10:05pm UTC
func isWeekend(t time.Time) bool {
	t = t.UTC()
	switch t.Weekday() {
	case time.Friday:
		h, _, _ := t.Clock()
		if h >= 12+10 {
			return true
		}
	case time.Saturday:
		return true
	case time.Sunday:
		h, m, _ := t.Clock()
		if h < 12+10 {
			return true
		}
		if h == 12+10 && m <= 5 {
			return true
		}
	}
	return false
}

//Tick get called when there is new data coming in
func (mrebstrat *MonthlyStrat) PerBar(k int, callback K.Callback) error {

	//assetCount := 0
	/* 	/* for asset := range mrebstrat.roc120 {
		if strings.ToLower(asset.Name) == "omx" {
			if mrebstrat.close[asset][k-asset.AdjK] > mrebstrat.up[asset][k-asset.AdjK-1] {
				mrebstrat.bull = true
			} else if mrebstrat.close[asset][k-asset.AdjK] < mrebstrat.dn[asset][k-asset.AdjK-1] {
				mrebstrat.bull = false
			}
		}
	} */
	/*assetCount++
	if mrebstrat.close[asset][k-asset.AdjK] > mrebstrat.ma200[asset][k-asset.AdjK] {
		mrebstrat.trend++
	}*/

	if callback.Date().Day() == 28 || callback.Date().Day() == 29 || callback.Date().Day() == 30 {
		//mrebstrat.oldMonth = callback.Date().Month()

		for asset := range mrebstrat.up {
			//fmt.Printf("Close %f mrebstrat %f\n", mrebstrat.close[asset][k-asset.AdjK], mrebstrat.up[asset][k-asset.AdjK])

			//Sell here
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
						"Cash": fmt.Sprintf("%f", cash),
					}})
			}

			close := mrebstrat.close[asset][k-asset.AdjK]

			/*keys66 := make([]roc, 0, len(mrebstrat.roc60))
			for asset, values := range mrebstrat.roc60 {
				if k-asset.AdjK > 0 {
					if k-asset.AdjK < len(values)-1 {
						keys66 = append(keys66, roc{asset.Name, asset, values[k-asset.AdjK] * mrebstrat.roc120[asset][k-asset.AdjK]})
					}
				}
			}*/

			/*sort.Slice(keys66, func(i, j int) bool {
				return keys66[i].value > keys66[j].value
			})*/

			if mrebstrat.roc60[asset][k-asset.AdjK-1] > 20 &&
				mrebstrat.roc60[asset][k-asset.AdjK-1] < 100 &&
				mrebstrat.roc120[asset][k-asset.AdjK] > 50 &&
				mrebstrat.roc120[asset][k-asset.AdjK] < 300 {
				//mrebstrat.bull {

				/*for i := 0; i < 5; i++ {

				if keys66[i].value > 5000 { */

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
							"ROC60":  fmt.Sprintf("%f", mrebstrat.roc60[asset][k-1]),
							"ROC120": fmt.Sprintf("%f", mrebstrat.roc120[asset][k-1]),
							"Cost":   fmt.Sprintf("%f", cost),
						}})
				}
			}
		}
	}

	return nil
}

func (mrebstrat *MonthlyStrat) End(k int, callback K.Callback) error {
	return nil
}

//Update is called before perBar
func (mrebstrat *MonthlyStrat) Update(k *int) error { return nil }

//OrderEvent gets called on order events
func (mrebstrat *MonthlyStrat) OrderEvent(orderEvent K.Event) {}
