package indicators

import (
	log "github.com/sirupsen/logrus"
)

type Series struct {
	*IndicatorStruct
}

func TimeSeries(series []float64) *Series {
	ts := &Series{
		IndicatorStruct: &IndicatorStruct{},
	}
	ts.IndicatorStruct.defaultValues = series

	log.WithFields(log.Fields{
		"size of series": len(series),
		"size of indicator struct default values": len(ts.IndicatorStruct.defaultValues),
	}).Debug("TIMESERIES> CREATED")

	return ts
}
