package generationk

import "math"

// Sum returns the sum of all elements of 'data'.
func Sum(data []float64) float64 {

	var sum float64

	for _, value := range data {
		sum += value
	}

	return sum
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
