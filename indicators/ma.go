package indicators

//NewSimpleMovingAverage is to start a new moving average
func SimpleMovingAverage(series []float64, period int) []float64 {
	returnSeries := make([]float64, period)

	for i := period; i < len(series); i++ {
		var sum float64
		for k := i; k < i+period; k++ {
			sum += series[k]
		}
		returnSeries[i] = sum / float64(period)
	}

	return returnSeries
}
