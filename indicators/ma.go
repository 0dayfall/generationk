package generationk

//MA is the Moving Average, other alias are SMA or SimpleMovingAverage
func MA(series []float64, period int) []float64 {
	return SMA(series, period)
}

//SMA is the Moving Average, other alias are MA or SimpleMovingAverage
func SMA(series []float64, period int) []float64 {
	return SimpleMovingAverage(series, period)
}

//SimpleMovingAverage to start a new moving average
func SimpleMovingAverage(series []float64, period int) []float64 {
	seriesLength := len(series)
	returnSeries := make([]float64, seriesLength)

	for i := 0; i < seriesLength; i++ {
		var sum float64

		for k := i - period; k < i && k > 0; k++ {
			sum += series[k]
		}
		/*
			for i := period; i <= len(slice); i++ {
				smaSlice = append(smaSlice, Sum(slice[i-period:i])/float64(period))
			}
		*/
		returnSeries[i] = sum / float64(period)
	}

	return returnSeries
}

//EMA is to start a new exponential moving average
//EMA = Price(t) * k + EMA(y) * (1 – k)
//t = today, y = yesterday, N = number of days in EMA, k = 2/(N+1)
func EMA(series []float64, period int) []float64 {
	return ExponentialMovingAverage(series, period)
}

//ExponentialMovingAverage is to start a new exponential moving average
func ExponentialMovingAverage(series []float64, period int) []float64 {
	//float k = 2 / (numberOfDays + 1);
	//return todaysPrice * k + EMAYesterday * (1 – k);

	seriesLength := len(series)
	returnSeries := make([]float64, seriesLength)

	decay := 2.0 / (float64(period) + 1)

	returnSeries[0] = (series[0] * decay) + (series[0] * (1.0 - decay))
	for i := 1; i < seriesLength; i++ {
		returnSeries[i] = (series[i] * decay) + (returnSeries[i-1] * (1.0 - decay))
	}

	return returnSeries
}
