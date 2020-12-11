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
