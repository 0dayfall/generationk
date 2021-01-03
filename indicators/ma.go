package indicators

//SimpleMovingAverage is the simple moving average
type SimpleMovingAverage struct {
	*IndicatorStruct
	dataType OHLC
}

//NewSimpleMovingAverage is to start a new moving average
func NewSimpleMovingAverage(value OHLC, period int) SimpleMovingAverage {
	ma := SimpleMovingAverage{
		IndicatorStruct: &IndicatorStruct{name: "Simple Moving Average", period: period, values: []float64{}},
		dataType:        value,
	}

	return ma
}

func NewSimpleMovingAverageParams(value OHLC, params Param) SimpleMovingAverage {
	ma := SimpleMovingAverage{
		IndicatorStruct: &IndicatorStruct{name: "Simple Moving Average", values: []float64{}, params: params},
		dataType:        value,
	}
	return ma
}

//Update is used to give data to the indicator
func (sma *SimpleMovingAverage) Update(values []float64) {
	sma.IndicatorStruct.values = values
}

//GetDataType is used to know which float64 array to use on the indicator
func (sma SimpleMovingAverage) GetDataType() OHLC {
	return sma.dataType
}

//Return the value at ix
func (sma *SimpleMovingAverage) Historic(ix int) float64 {
	if len((*sma.IndicatorStruct).values) < 1 {
		return 0.0
	}
	var sum float64
	for k := 0; k <= sma.period-1; k++ {
		sum += sma.IndicatorStruct.values[k]
	}
	returnValue := sum / float64(sma.period)

	return returnValue
}
