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

	records[1] = stripThousandSeparator(records[1])
	floats[0], err = strconv.ParseFloat(records[1], 64)
	if err != nil {
		return time, nil, err
	}

	records[2] = stripThousandSeparator(records[2])
	floats[1], err = strconv.ParseFloat(records[2], 64)
	if err != nil {
		return time, nil, err
	}

	records[3] = stripThousandSeparator(records[3])
	floats[2], err = strconv.ParseFloat(records[3], 64)
	if err != nil {
		return time, nil, err
	}

	records[4] = stripThousandSeparator(records[4])
	floats[3], err = strconv.ParseFloat(records[4], 64)
	if err != nil {
		return time, nil, err
	}

	records[5] = stripThousandSeparator(records[5])
	floats[4], err = strconv.ParseFloat(records[5], 64)
	if err != nil {
		return time, nil, err
	}

	return time, floats, nil
}
