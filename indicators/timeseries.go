package indicators

import (
	log "github.com/sirupsen/logrus"
)

//TimeSeries structure is used to create a time series which may be updated
type TimeSeries struct {
	*IndicatorStruct
	dataType OHLC
}

//NewTimeSeries is used to create a time series which may be updates
func NewTimeSeries(value OHLC, period int) TimeSeries {

	ts := TimeSeries{
		IndicatorStruct: &IndicatorStruct{
			name:   "Time Series",
			period: period},
		dataType: value,
	}

	log.WithFields(log.Fields{
		"type of value in time series": value,
	}).Debug("TIMESERIES> CREATED")

	return ts
}

//Update is used to initialize indicator with values
func (ts TimeSeries) Update(values []float64) {
	log.WithFields(log.Fields{
		"values length": len(values),
	}).Debug("TIMESERIES> Update()")
	ts.IndicatorStruct.values = values
}

//GetDataType explains to generationk which type it should be updated with
func (ts TimeSeries) GetDataType() OHLC {
	return ts.dataType
}
