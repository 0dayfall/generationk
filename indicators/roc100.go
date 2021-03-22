package generationk

//ROC100 is the Rate Of Change, other alias are RateOfChange
func ROC100(series []float64, period int) []float64 {
	return RateOfChange100(series, period)
}

//RateOfChange100 is used to get a value between 1 and 100 to pin point the rate of change
//ROC=(previous value / current value−1)∗100
func RateOfChange100(series []float64, period int) []float64 {
	seriesLength := len(series)
	returnSeries := make([]float64, seriesLength)

	//Adjust to the array starting with index 0
	period -= 1

	for i := 0; i < seriesLength; i++ {
		//series[68] - series[68-66] / series[68-66]
		if i-period < 0 {
			returnSeries[i] = 0
		} else {
			returnSeries[i] = 100 * (series[i] - series[i-period]) / series[i-period]
		}
	}

	return returnSeries
}
