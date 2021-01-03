package indicators

//TimeSeries structure is used to create a time series which may be updated,
//Should be used with open, high, low, close or volume
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

	return ts
}

//Update is used to initialize indicator with values
func (ts TimeSeries) Update(values []float64) {
	ts.IndicatorStruct.values = values
}

//GetDataType explains to generationk which type it should be updated with
func (ts TimeSeries) GetDataType() OHLC {
	return ts.dataType
}
