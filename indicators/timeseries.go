package indicators

import (
	log "github.com/sirupsen/logrus"
)

//TimeSeries structure is used to create a time series which may be updated
type TimeSeries struct {
	*IndicatorStruct
}

//NewTimeSeries is used to create a time series which may be updates
func NewTimeSeries(series []float64) (*TimeSeries, error) {
	if len(series) < 1 {
		return nil, IndicatorNotReadyError{
			Msg: "NewTimeSeries",
			Len: len(series),
		}
	}
	ts := &TimeSeries{
		IndicatorStruct: &IndicatorStruct{},
	}
	ts.IndicatorStruct.values = series

	log.WithFields(log.Fields{
		"size of series": len(series),
		"size of indicator struct default values": len(ts.IndicatorStruct.values),
	}).Debug("TIMESERIES> CREATED")

	return ts, nil
}
