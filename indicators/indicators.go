package indicator

//I type interface
type I interface {
	LargerThan(number float64, number2 float64)
	Value(value string)
}

//MA struct
type indicator struct {
}

//LargerThan the larger than operator
func (m *indicator) LargerThan(number float64) {}

//Value the default value of the indicator
func (m *indicator) Value(value string) {}
