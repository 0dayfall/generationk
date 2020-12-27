package indicators

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

//OHLC is the value which will be used for the indicator
type OHLC int

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
	GetDataType() OHLC
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

func (m IndicatorStruct) GetName() string {
	return m.name
}

func (m IndicatorStruct) GetPeriod() int {
	return m.period
}

func (m IndicatorStruct) Value() float64 {
	return m.values[0]
}

func (m *IndicatorStruct) ValueAtIndex(index int) float64 {
	if len((*m).values) < 1 {
		log.Error("GENERIC> No values in generic indicator")
		return 0.0
	}
	log.WithFields(log.Fields{
		"index":                  index,
		"len":                    len(m.values),
		"m.defaultValues[index]": m.values[index],
	}).Debug("GENERIC INDICATOR>ValueAtIndex()")
	return m.values[index]
}

func (m IndicatorStruct) Values() []float64 {
	return m.values
}

//LargerThan the larger than operator
func (m IndicatorStruct) LargerThanValue(ind Indicator) bool {
	fmt.Printf("Larger than Value is used\n")
	return m.Value() > m.Value()
}

//LargerThanInd
func (m IndicatorStruct) LargerThanInd(ind Indicator) []bool {
	result := make([]bool, len(m.Values()))
	for i, j := range m.Values() {
		result[i] = m.ValueAtIndex(i) > j
	}
	return result
}
