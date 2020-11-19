package indicators

// MACD stands for moving average convergence divergence.
func MACD(data mfloat, ema ...int) ([]float64, []float64) {

	var macd, ema1, ema2, ema3 mfloat

	if len(ema) < 3 {
		ema = []int{12, 26, 9}
	}

	ema1 = data.EMA(ema[0])
	ema2 = data.EMA(ema[1])
	macd = SubSlices(ema1, ema2)
	ema3 = macd.EMA(ema[2])

	return macd, ema3
}
