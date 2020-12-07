package internal

import (
	"encoding/csv"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//CSVDataManager type
type CSVDataManager struct {
	dataChannel chan Event
}

func (d *CSVDataManager) getData(period int) []float64 {
	return nil
}

func (d *CSVDataManager) getLatestData() float64 {
	return 0.0
}

func parseFloat(value string) float64 {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Was not possible to parse float: %s", value)
		return 0.0
	}
	return floatValue
}

//NewCSVDataManager creates a new data manager object
func NewCSVDataManager(ctx *Context) CSVDataManager {
	dm := CSVDataManager{
		dataChannel: ctx.EventChannel(),
	}
	log.WithFields(log.Fields{
		"dataChannel": dm.dataChannel,
	}).Debug("Created CSVDataManager")
	return dm
}

//ReadCSVFile reads a CSV file
func (d *CSVDataManager) ReadCSVFile(file string) {
	name := strings.TrimSuffix(filepath.Base(file), path.Ext(file))
	csvfile, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"File name": file,
	}).Debug("CSVDataManager> OPENED FILE")

	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()

	if err != nil && err != io.EOF {
		log.Error("Was not possible to read the file %s", err)
	}

	s := make([]OHLC, len(records))
	for i, record := range records {
		// Read each record from csv
		record1, err := time.Parse("1/2/2006 00:00:00", record[0]+" "+record[1])
		record2, err := strconv.ParseFloat(record[2], 64)
		record3, err := strconv.ParseFloat(record[3], 64)
		record4, err := strconv.ParseFloat(record[4], 64)
		record5, err := strconv.ParseFloat(record[5], 64)
		record6, err := strconv.Atoi(record[6])

		if err != nil {
			log.Error("Was not possible to parse the format on  line %d, %s", i, err)
		}

		ohlc := OHLC{Time: record1, Open: record2, High: record3, Low: record4, Close: record5, Volume: record6}
		log.WithFields(log.Fields{
			"Name": name,
			"Ohlc": ohlc,
		}).Debug("DataEvent$ ")
		d.dataChannel <- DataEvent{Name: name, Ohlc: ohlc}
		s[i] = ohlc
	}

}
