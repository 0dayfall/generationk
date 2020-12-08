package indicators

import (
	log "github.com/sirupsen/logrus"
)

type Series struct {
	*IndicatorStruct
}

func NewTimeSeries(series []float64) (*Series, error) {
	if len(series) < 1 {
		return nil, IndicatorNotReadyError{
			Msg: "SimpleMovingAverage",
			Len: len(series),
		}
	}
	ts := &Series{
		IndicatorStruct: &IndicatorStruct{},
	}
	ts.IndicatorStruct.defaultValues = series

	log.WithFields(log.Fields{
		"size of series": len(series),
		"size of indicator struct default values": len(ts.IndicatorStruct.defaultValues),
	}).Debug("TIMESERIES> CREATED")

	return ts, nil
}
