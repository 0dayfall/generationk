package indicators

import (
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"
)

type mfloat []float64

//OhlcValue is the value which will be used for the indicator
type OhlcValue int

const (
	//Open value will be used for the data
	Open = iota
	//High value will be used for the data
	High
	//Low value will be used for the data
	Low
	//Close value will be used for the data
	Close
	//Volume value will be used for the data
	Volume
	//Default is used when the data is something else than ohlcv
	Default
)

//Updatable tells generatinok this it is possible to run Update()
type Updatable interface {
	Update()
}

//Comparison type interface
type Comparison interface {
	LargerThanValue(number2 float64) bool
	LargerThanInd(ind Indicator)
}

//Indicator Interface for all indicators
type Indicator interface {
	ValueAtIndex(index int) float64
	Value() float64
	Values() []float64
	Update(values []float64)
	GetDataType() OhlcValue
	GetName() string
	GetPeriod() int
}

//IndicatorStruct contains a default set of values
type IndicatorStruct struct {
	//ctx *genk.Context
	name   string
	period int
	values []float64
}

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

func (m IndicatorStruct) GetPeriod() int {
	return m.period
}

func (m IndicatorStruct) GetName() string {
	return m.name
}

func (m *IndicatorStruct) Shift() {
	if len(m.values) > 1 {
		m.values = m.values[1:]
	}
}

func (m *IndicatorStruct) Value() float64 {
	return m.values[0]
}

func (m *IndicatorStruct) ValueAtIndex(index int) float64 {
	log.WithFields(log.Fields{
		"index":                  index,
		"len":                    len(m.values),
		"m.defaultValues[index]": m.values[index],
	}).Debug("GENERIC INDICATOR>ValueAtIndex()")
	return m.values[index]
}

func (m *IndicatorStruct) Values() []float64 {
	return m.values
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
