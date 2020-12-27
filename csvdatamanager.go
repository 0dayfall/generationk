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
	//dataChannel chan Event
	callback DataHandler
}

/*func (d *CSVDataManager) getData(period int) []float64 {
	return nil
}

func (d *CSVDataManager) getLatestData() float64 {
	return 0.0
}*/

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
	/*log.WithFields(log.Fields{
		"dataChannel": dm.dataChannel,
	}).Debug("Created CSVDataManager")*/
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

//ReadFolderWithCSVFilesAsync is used to read a folder of files and put them on the queue to the strategy
func (d CSVDataManager) ReadFolderWithCSVFilesAsync(folder string, wg *sync.WaitGroup) {
	//var heap OhlcHeap
	files, err := filepath.Glob(folder + "*.csv")
	if err != nil {
		log.Fatal(err)
	}

	d.ReadCSVFilesAsync(files, wg)
}

//ReadCSVFile reads a CSV file
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

func (d *CSVDataManager) callbackOnDataEvent(name string, ohlc []OHLC) int {
	var count int

	for k := range ohlc {
		d.callback.DataEvent(DataEvent{Name: name, Ohlc: ohlc[k]})
		count++
		//fmt.Printf("%d", k)
	}

	return count
	//d.dataChannel <- Quit{}
}

/*func (d *CSVDataManager) putDataOnChannel(name string, ohlc []OHLC) {

	for k := 0; k < len(ohlc); k++ {
		d.dataChannel <- DataEvent{Name: name, Ohlc: ohlc[k]}
	}
	d.dataChannel <- Quit{}
}*/

func (d *CSVDataManager) read(file string) int {
	ohlc := d.readCSVFile(file)
	name := strings.TrimSuffix(filepath.Base(file), path.Ext(file))
	//d.putDataOnChannel(name, ohlc)
	return d.callbackOnDataEvent(name, ohlc)
}

//ReadCSVFile is used to start a go thread
func (d *CSVDataManager) ReadCSVFile(file string) int {
	return d.read(file)
}

//ReadCSVFileAsync is used to start a go thread
func (d *CSVDataManager) ReadCSVFileAsync(file string, wg *sync.WaitGroup) {
	d.read(file)
	wg.Done()
}
