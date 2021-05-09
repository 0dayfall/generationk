package generationk

import "math"

// The Formula for Donchian Channels Is:
// UC = Highest High in Last N Periods
// Middle Channel=((UC−LC)/2)
// LC = Lowest Low in Last N periods
// where:
// UC = Upper channel
// N = Number of minutes, hours, days, weeks, months
// Period = Minutes, hours, days, weeks, months
// LC = Lower channel
func Bollinger(series []float64, mult int, period int) ([]float64, []float64, []float64) {
	seriesLength := len(series)

	up := make([]float64, seriesLength)
	middle := make([]float64, seriesLength)
	dn := make([]float64, seriesLength)

	for i := 0; i < seriesLength; i++ {
		var sum float64

		for k := i - period; k < i && k >= 0; k++ {
			sum += series[k]
		}

		middle[i] = sum / float64(period)
	}

	var stdDev float64

	for j := 0; j < seriesLength; j++ {

		for k := j - period; k < j && k >= 0; k++ {
			// The use of Pow math function func Pow(x, y float64) float64
			stdDev += math.Pow(series[k]-middle[k], 2)
		}

		// The use of Sqrt math function func Sqrt(x float64) float64
		stdDev = math.Sqrt(stdDev / float64(period))
		up[j] += stdDev * float64(mult)
		dn[j] -= stdDev * float64(mult)

	}

	return up, middle, dn
}
