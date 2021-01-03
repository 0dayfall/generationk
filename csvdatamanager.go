package generationk

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

//CSVDataManager type is used to send DataEvents via callback to generationK
type CSVDataManager struct {
	callback DataHandler
}

//pasetFloat is used to parse the floats from the CSV files and is a better way to
//to handle errors
func parseFloat(value string) float64 {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal(err)
	}

	return floatValue
}

//NewCSVDataManager creates a new data manager object
func NewCSVDataManager(dataHandler DataHandler) *CSVDataManager {

	dm := &CSVDataManager{
		//dataChannel: genk.market.eventChannel,
		callback: dataHandler,
	}

	return dm
}

//ReadCSVFilesAsync is used to read files asynchronous
func (d CSVDataManager) ReadCSVFilesAsync(files []string, wg *sync.WaitGroup) {
	for k := 0; k < len(files); k++ {
		wg.Add(1)
		go d.ReadCSVFileAsync(files[k], wg)
	}

	wg.Wait()
}

//ReadFolderWithCSVFilesAsync is used to read a folder of files
func (d CSVDataManager) ReadFolderWithCSVFilesAsync(folder string, wg *sync.WaitGroup) {
	//var heap OhlcHeap
	files, err := filepath.Glob(folder + "*.csv")
	if err != nil {
		log.Fatal(err)
	}

	d.ReadCSVFilesAsync(files, wg)
}

//ReadCSVFile reads a CSV file and maps the records according to this method
func (d CSVDataManager) readCSVFile(file string) []OHLC {

	csvfile, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}

	defer csvfile.Close()

	log.WithFields(log.Fields{
		"File name": file,
	}).Debug("CSVDataManager> OPENED FILE")

	// Parse the file
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()

	if err != nil && errors.Is(err, io.EOF) {
		log.Fatal(err)
	}

	s := make([]OHLC, len(records))

	for i, record := range records {
		// Read each record from csv
		record1, err := time.Parse("1/2/2006 00:00:00", record[0]+" "+record[1])
		if err != nil {
			log.Fatal(err)
		}

		record2 := parseFloat(record[2])
		record3 := parseFloat(record[3])
		record4 := parseFloat(record[4])
		record5 := parseFloat(record[5])

		record6, err := strconv.Atoi(record[6])

		if err != nil {
			log.Fatal(err)
		}

		ohlc := OHLC{time: record1, open: record2, high: record3, low: record4, close: record5, volume: record6}
		s[i] = ohlc
	}

	return s
}

//callbackOnDataEvent is used to send each row read from the CSV file to the callback
func (d *CSVDataManager) callbackOnDataEvent(name string, ohlc []OHLC) int {
	var count int

	for k := range ohlc {
		d.callback.DataEvent(DataEvent{Name: name, Ohlc: ohlc[k]})
		count++
	}

	return count
}

//read is used to read a single file and feed back the data to the callback
func (d *CSVDataManager) read(file string) int {
	ohlc := d.readCSVFile(file)
	name := strings.TrimSuffix(filepath.Base(file), path.Ext(file))
	return d.callbackOnDataEvent(name, ohlc)
}

//ReadCSVFile is used to read a single file
func (d *CSVDataManager) ReadCSVFile(file string) int {
	return d.read(file)
}

//ReadCSVFileAsync is a sigle file asynchronously
func (d *CSVDataManager) ReadCSVFileAsync(file string, wg *sync.WaitGroup) {
	d.read(file)
	wg.Done()
}
