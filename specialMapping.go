package generationk

import (
	"log"
	"time"
)

func MapRecordsInvesting(records ...string) (time.Time, []float64) {
	//time, err := time.Parse("1/2/2006 15:04:05", records[0]+" "+records[1])
	time, err := time.Parse("2006-01-02 15:04:05", records[0])
	if err != nil {
		log.Fatal(err)
	}

	floats := make([]float64, 5)

	floats[0] = ParseFloat(records[1])
	floats[1] = ParseFloat(records[2])
	floats[2] = ParseFloat(records[3])
	floats[3] = ParseFloat(records[4])
	floats[4] = ParseFloat(records[5])

	return time, floats
}
