package generationk

import (
	"strconv"
	"time"
)

func MapRecordsInvesting(records ...string) (time.Time, []float64, error) {

	time, err := time.Parse("2006-01-02 15:04:05", records[0])
	if err != nil {
		return time, nil, err
	}

	floats := make([]float64, 5)

	floats[0], err = strconv.ParseFloat(records[1], 64)
	if err != nil {
		return time, nil, err
	}

	floats[1], err = strconv.ParseFloat(records[2], 64)
	if err != nil {
		return time, nil, err
	}

	floats[2], err = strconv.ParseFloat(records[3], 64)
	if err != nil {
		return time, nil, err
	}

	floats[3], err = strconv.ParseFloat(records[4], 64)
	if err != nil {
		return time, nil, err
	}

	floats[4], err = strconv.ParseFloat(records[5], 64)
	if err != nil {
		return time, nil, err
	}

	return time, floats, nil
}
