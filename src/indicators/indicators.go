package indicators

//I type interface
type ComparisonInterface interface {
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
type Comparison struct {
	values []float64
}

func (m *Comparison) Value() float64 {
	return m.values[0]
}

func (m *Comparison) ValueIndex(index int) float64 {
	return m.values[index]
}

func (m *Comparison) Values() []float64 {
	return m.values
}

//Value the default value of the indicator
func (m *Comparison) ValueAttr(value string) {}

//LargerThan the larger than operator
func (m *Comparison) LargerThanValue(ind Indicator) bool {
	return m.Value() > ind.Value()
}

//LargerThanInd
func (m *Comparison) LargerThanInd(ind Indicator) []bool {
	result := make([]bool, len(ind.Values()))
	for i, j := range ind.Values() {
		result[i] = m.values[i] > j
	}
	return result
}
