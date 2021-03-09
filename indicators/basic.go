package indicators

import "math"

// Sum returns the sum of all elements of 'data'.
func Sum(data []float64) float64 {

	var sum float64

	for _, value := range data {
		sum += value
	}

	return sum
}

//This is not used
type mfloat []float64

// Ema calculates exponential moving average of a slice for a certain
// number of tiSmame periods.
func (slice mfloat) EMA(period int) []float64 {

	var emaSlice []float64

	ak := period + 1
	k := float64(2) / float64(ak)

	emaSlice = append(emaSlice, slice[0])

	for i := 1; i < len(slice); i++ {
		emaSlice = append(emaSlice, (slice[i]*float64(k))+(emaSlice[i-1]*float64(1-k)))
	}

	return emaSlice
}

// Sma calculates simple moving average of a slice for a certain
// number of time periods.
func (slice mfloat) SMA(period int) []float64 {

	var smaSlice []float64

	for i := period; i <= len(slice); i++ {
		smaSlice = append(smaSlice, Sum(slice[i-period:i])/float64(period))
	}

	return smaSlice
}

// Avg returns 'data' average.
func Avg(data []float64) float64 {

	return Sum(data) / float64(len(data))
}

func Std(slice []float64, period int) float64 {
	return StandardDeviation(slice, period)
}

// Std returns standard deviation of a slice.
func StandardDeviation(series []float64, period int) float64 {

	lengthSeries := len(series)
	result := make([]float64, lengthSeries)

	//The average is the starting point
	ma := MA(series, period)

	//The difference between each data point and the average is calculated and then the values are ^2
	for i := 0; i < lengthSeries; i++ {
		result[i] = math.Pow(series[i]-ma[i], 2)
	}

	/*TODO: FINISH THIS*/
	for i := 0; i < lengthSeries; i++ {
		result[i] = math.Sqrt(Sum(result) / float64(period))
	}
	//The variance is the average of these values and the standard deviation is the square root of the variance
	return math.Sqrt(Sum(result) / float64(lengthSeries))
}

// AddToAll adds a value to all slice elements.
func AddToAll(slice []float64, val float64) []float64 {

	var addedSlice []float64

	for i := 0; i < len(slice); i++ {
		addedSlice = append(addedSlice, slice[i]+val)
	}

	return addedSlice
}

// SubSlices subtracts two slices.
func SubSlices(slice1, slice2 []float64) []float64 {

	var result []float64

	for i := 0; i < len(slice1); i++ {
		result = append(result, slice1[i]-slice2[i])
	}

	return result
}

// AddSlices adds two slices.
func AddSlices(slice1, slice2 []float64) []float64 {

	var result []float64

	for i := 0; i < len(slice1); i++ {
		result = append(result, slice1[i]+slice2[i])
	}

	return result
}

// DivSlice divides a slice by a float.
func DivSlice(slice []float64, n float64) []float64 {

	var result []float64

	for i := 0; i < len(slice); i++ {
		result = append(result, slice[i]/n)
	}

	return result
}
