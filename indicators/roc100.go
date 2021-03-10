package indicators

import "fmt"

//ROC100 is the Rate Of Change, other alias are RateOfChange
func ROC100(series []float64, period int) []float64 {
	return RateOfChange100(series, period)
}

//RateOfChange100 is used to get a value between 1 and 100 to pin point the rate of change
//ROC=(previous value / current value−1)∗100
func RateOfChange100(series []float64, period int) []float64 {
	seriesLength := len(series)
	returnSeries := make([]float64, seriesLength)

	for i := period; i < seriesLength; i++ {
		returnSeries[period] = (series[i] - series[i-period]) / series[i-period] * 100
	}
	fmt.Println("Return the rate of change 100")
	return returnSeries
}
