package indicators

import log "github.com/sirupsen/logrus"

//SimpleMovingAverage for the basics
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

	log.WithFields(log.Fields{
		"type of value in moving average series": value,
	}).Debug("MA> CREATED")

	return ma
}

//Update is used to give data to the indicator
func (sma *SimpleMovingAverage) Update(values []float64) {
	log.WithFields(log.Fields{
		"values length": len(values),
	}).Debug("MA> Update()")
	sma.IndicatorStruct.values = values
	log.WithFields(log.Fields{
		"sma.IndicatorStruct.values": len(sma.IndicatorStruct.values),
	}).Debug("MA> Update()")
}

//GetDataType is used to know which float64 array to use on the indicator
func (sma SimpleMovingAverage) GetDataType() OHLC {
	return sma.dataType
}

func (sma *SimpleMovingAverage) ValueAtIndex(ix int) float64 {
	if len((*sma.IndicatorStruct).values) < 1 {
		log.Error("No values in SimpleMovingAverage")
		return 0.0
	}
	var sum float64
	for k := 0; k <= sma.period-1; k++ {
		sum += sma.IndicatorStruct.values[k]
	}
	returnValue := sum / float64(sma.period)
	log.WithFields(log.Fields{
		"sma.IndicatorStruct.values[ix]": returnValue,
	}).Debug("MA> ValueAtIndex()")
	return returnValue
}
