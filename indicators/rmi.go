package generationk

import (
	"math"
)

/************************** CALCULATION *****************************
FastemaInc      = ema(max(close - close[FastMomentum], 0), FastLenght)
FastemaDec      = ema(max(close[FastMomentum] - close, 0), FastLenght)
FastRMI         = FastemaDec == 0 ? 0 : 50 - 100 / (1 + FastemaInc / FastemaDec)

SlowemaInc      = ema(max(close - close[SlowMomentum], 0), SlowLenght)
SlowemaDec      = ema(max(close[SlowMomentum] - close, 0), SlowLenght)
SlowRMI         = SlowemaDec == 0 ? 0 : 50 - 100 / (1 + SlowemaInc / SlowemaDec)

*********************************************************************/
func RMI(series []float64, fastRMIlength, fastMomentum, slowRMIlength, slowMomentum int) ([]float64, []float64) {

	fast := rmiHelper(series, fastRMIlength, fastMomentum)

	slow := rmiHelper(series, slowRMIlength, slowMomentum)

	return fast, slow
}

func rmiHelper(series []float64, RMIlength, momentum int) []float64 {
	size := len(series)
	results := make([]float64, size)
	emaInc := make([]float64, size)
	emaDec := make([]float64, size)

	var lookback int
	for i := range series {

		lookback = i - momentum
		if lookback < 0 {
			lookback = 0
		}

		emaInc[i] = math.Max(series[i]-series[lookback], 0.0)
		emaDec[i] = math.Max(series[lookback]-series[i], 0.0)

	}
	emaInc = EMA(emaInc, RMIlength)
	emaDec = EMA(emaDec, RMIlength)

	for i := range results {
		if emaDec[i] == 0 {
			results[i] = 0
		} else {
			results[i] = 50 - 100/(1+emaInc[i]/emaDec[i])
			/*fmt.Printf("emaInc value: %f ", emaInc[i])
			fmt.Printf("emaDec value: %f ", emaDec[i])
			fmt.Printf("RMI value: %f\n", results[i])*/
		}
	}

	return results
}
