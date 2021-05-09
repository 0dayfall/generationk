package generationk

import (
	"math"
)

// The Formula for Donchian Channels Is:
// UC = Highest High in Last N Periods
// Middle Channel=((UC−LC)/2)
// LC = Lowest Low in Last N periods
// where:
// UC = Upper channel
// N = Number of minutes, hours, days, weeks, months
// Period = Minutes, hours, days, weeks, months
// LC = Lower channel
func Donchian(high []float64, low []float64, period int) ([]float64, []float64, []float64) {
	highLength := len(high)
	lowLength := len(low)

	if highLength != lowLength {
		//return 0,0,0
	}

	up := make([]float64, highLength)
	middle := make([]float64, highLength)
	down := make([]float64, lowLength)

	for i := 0; i < highLength; i++ {

		max, mid := 0.0, 0.0
		min := 10000000.0

		var window int
		if i-period < 0 {
			window = 0
		} else {
			window = i - period + 1
		}

		for k := window; k <= i; k++ {

			max = math.Max(max, high[k])
			min = math.Min(min, low[k])

			//Need the max and min to calculate the mid
			mid = (max + min) / 2

		}

		up[i] = max
		middle[i] = mid
		down[i] = min

	}

	return up, middle, down
}
