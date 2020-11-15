package indicators

import (
	"fmt"
	"math"
)

type mfloat []float64

//I type interface
type Comparison interface {
	LargerThanValue(number float64, number2 float64)
	LargerThanInd(ind Indicator)
}

//Indicator Interface for all indicators
type Indicator interface {
	ValueAtIndex(index int) float64
	Value() float64
	Values() []float64
}

//Indicator struct
type IndicatorStruct struct {
	defaultValues []float64
}

// Avg returns 'data' average.
func Avg(data []float64) float64 {

	return Sum(data) / float64(len(data))
}

// Sum returns the sum of all elements of 'data'.
func Sum(data []float64) float64 {

	var sum float64

	for _, x := range data {
		sum += x
	}

	return sum
}

// Std returns standard deviation of a slice.
func Std(slice []float64) float64 {

	var result []float64

	mean := Avg(slice)

	for i := 0; i < len(slice); i++ {
		result = append(result, math.Pow(slice[i]-mean, 2))
	}

	return math.Sqrt(Sum(result) / float64(len(result)))
}

// AddToAll adds a value to all slice elements.
func (slice mfloat) AddToAll(val float64) []float64 {

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

func (m *IndicatorStruct) Value() float64 {
	return m.defaultValues[0]
}

func (m *IndicatorStruct) ValueAtIndex(index int) float64 {
	return m.defaultValues[index]
}

func (m *IndicatorStruct) Values() []float64 {
	return m.defaultValues
}

//Value the default value of the indicator
func (m *IndicatorStruct) ValueAttr(value string) {}

//LargerThan the larger than operator
func (m *IndicatorStruct) LargerThanValue(ind Indicator) bool {
	fmt.Printf("Larger than Value is used\n")
	return m.Value() > ind.Value()
}

//LargerThanInd
func (m *IndicatorStruct) LargerThanInd(ind Indicator) []bool {
	result := make([]bool, len(ind.Values()))
	for i, j := range ind.Values() {
		result[i] = m.ValueAtIndex(i) > j
	}
	return result
}
