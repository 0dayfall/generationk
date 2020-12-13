package indicators

import (
	log "github.com/sirupsen/logrus"
)

//TimeSeries structure is used to create a time series which may be updated
type TimeSeries struct {
	*IndicatorStruct
	value OhlcValue
}

//GetDataType explains to generationk which type it should be updated with
func (ts TimeSeries) GetDataType() OhlcValue {
	return ts.value
}

//Update is used to initialize indicator with values
func (ts TimeSeries) Update(values []float64) {
	log.WithFields(log.Fields{
		"values length": len(values),
	}).Debug("TIMESERIES> Update()")
	ts.IndicatorStruct.values = values
}

//NewTimeSeries is used to create a time series which may be updates
func NewTimeSeries(value OhlcValue, period int) *TimeSeries {

	ts := &TimeSeries{
		IndicatorStruct: &IndicatorStruct{
			name:   "Time Series",
			period: period},
		value: value,
	}

	log.WithFields(log.Fields{
		"type of value in time series": value,
	}).Debug("TIMESERIES> CREATED")

	return ts
}
