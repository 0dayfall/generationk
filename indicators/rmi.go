package indicators

import "math"

/************************** CALCULATION *****************************
FastemaInc      = ema(max(close - close[FastMomentum], 0), FastLenght)
FastemaDec      = ema(max(close[FastMomentum] - close, 0), FastLenght)
FastRMI         = FastemaDec == 0 ? 0 : 50 - 100 / (1 + FastemaInc / FastemaDec)

SlowemaInc      = ema(max(close - close[SlowMomentum], 0), SlowLenght)
SlowemaDec      = ema(max(close[SlowMomentum] - close, 0), SlowLenght)
SlowRMI         = SlowemaDec == 0 ? 0 : 50 - 100 / (1 + SlowemaInc / SlowemaDec)

*********************************************************************/
func Rmi(series []float64, fastMomentum int, fastRMI int, slowRmi int) ([]float64, []float64) {
	EMA(math.Max(series[i] - series[fastMomentum], 0))
}