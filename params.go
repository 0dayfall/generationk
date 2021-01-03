//+build dontbuild
package generationk

import indicators "github.com/0dayfall/generationk/indicators"

type Evaluation interface {
	Indicators() []indicators.Indicator
	EvaluateIndicator(indicators []indicators.Indicator)
	EvaluateData(ohlc []OHLC)
}

type MACross struct{}

func (ma MACross) Indicators() map[string]indicators.Indicator {
	var ind map[string]indicators.Indicator

	ma15 := indicators.NewSimpleMovingAverage(indicators.Close, 15)
	ma15.SetParams(15, 30)
	ind["ma15"] = &ma15

	return ind
}

func (ma MACross) EvaluateIndicator(indicators map[string]indicators.Indicator) bool {
	return indicators["ma15"].Current() > 15
}

func run() {
	ma15 := indicators.NewSimpleMovingAverage(indicators.Close, 15)
	ma15.SetParams(15, 30)

	//Or(indicators)
}

func Or(evaluations []Evaluation) bool {
	/*for _, eval := range evaluations {
		if eval == true {
			return true
		}
	}*/
	return false
}
