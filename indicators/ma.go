package indicators

//SimpleMovingAverage for the basics
type SimpleMovingAverage struct {
	IndicatorStruct
	period int
	//value  OhlcValue
}

//NewSimpleMovingAverage is to start a new moving average
func NewSimpleMovingAverage(value OhlcValue, period int) (Indicator, error) {
	ma := SimpleMovingAverage{
		IndicatorStruct: IndicatorStruct{},
		period:          period,
	}
	return &ma, nil
}

func (sma *SimpleMovingAverage) calculate(k int) {
	var sum float64
	for k := 0; k <= sma.period; k++ {
		sum += values[k]
	}
	sma.IndicatorStruct.values = sum
}
