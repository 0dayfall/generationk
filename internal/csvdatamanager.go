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

func (d CSVDataManager) readCSVFilesAsync(files []string) {
	for k := 0; k < len(files); k++ {
		d.ReadCSVFileAsync(files[k])
	}
}

//ReadFolderWithCSVFilesAsync is used to read a folder of files and put them on the queue to the strategy
func (d CSVDataManager) ReadFolderWithCSVFilesAsync(folder string) {
	//var heap OhlcHeap
	files, err := filepath.Glob("*.csv")
	if err != nil {
		log.Fatal(err)
	}
	d.readCSVFilesAsync(files)
}

//ReadCSVFile reads a CSV file
func (d CSVDataManager) readCSVFile(file string) []OHLC {

	csvfile, err := os.Open(file)
	defer csvfile.Close()

	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"File name": file,
	}).Debug("CSVDataManager> OPENED FILE")

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
			"Name of file": file,
			"Ohlc":         ohlc,
		}).Debug("DataEvent$ ")

		s[i] = ohlc
	}

	return s
}

func (d CSVDataManager) putDataOnChannel(name string, ohlc []OHLC) {

	for k := 0; k < len(ohlc); k++ {
		d.dataChannel <- DataEvent{Name: name, Ohlc: ohlc[k]}
	}
	d.dataChannel <- Quit{}
}

//ReadCSVFileAsync is used to start a go thread
func (d CSVDataManager) ReadCSVFileAsync(file string) {
	ohlc := d.readCSVFile(file)
	name := strings.TrimSuffix(filepath.Base(file), path.Ext(file))
	go d.putDataOnChannel(name, ohlc)
}
