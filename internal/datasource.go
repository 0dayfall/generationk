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

//DataManager type
type CSVDataManager struct {
	fieldMapper func(records [][]string) []OHLC
	dataChannel chan Event
	Asset       []Asset
	ctx         *Context
}

func (d *DataManager) getData(period int) []float64 {
	return nil
}

func (d *DataManager) getLatestData() float64 {
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

//NewDataManager creates a new data manager object
func NewDataManager(ctx *Context) DataManager {
	dm := DataManager{
		dataChannel: ctx.EventChannel(),
		Asset:       make([]Asset, 0),
	}
	log.WithFields(log.Fields{
		"dataChannel": dm.dataChannel,
	}).Debug("Created DataManager")
	return dm
}

func reverseSlice(ohlc *[]OHLC) {
	for i := len(*ohlc)/2 - 1; i >= 0; i-- {
		opp := len(*ohlc) - 1 - i
		(*ohlc)[i], (*ohlc)[opp] = (*ohlc)[opp], (*ohlc)[i]
	}
}

//ReadCSVFile reads a CSV file
func (d *DataManager) ReadCSVFile(file string) {
	asset := Asset{
		Name: strings.TrimSuffix(filepath.Base(file), path.Ext(file)),
	}
	log.WithFields(log.Fields{
		"Name": asset.Name,
	}).Debug("Asset created")
	//d.ctx.AddAsset(&asset)
	d.Asset = append(d.Asset, asset)
}

func (d *DataManager) Start() {
	csvfile, err := os.Open(d.Asset[0].Name)

	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"File name": d.Asset[0].Name,
	}).Debug("DATAMANAGER> OPENED FILE")

	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()

	if err != nil && err != io.EOF {
		log.Printf("Was not possible to read the file %s", err)
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
			log.Printf("Was not possible to parse the format on  line %d, %s", i, err)
		}

		ohlc := OHLC{Time: record1, Open: record2, High: record3, Low: record4, Close: record5, Volume: record6}
		d.dataChannel <- DataEvent{Name: d.Asset[0].Name, Ohlc: ohlc}
		s[i] = ohlc
		//fmt.Printf("In addValue: s is %v\n", s)
	}

}

/*func readCSV(universe []string) []Stock {
	stock := make([]Stock, len(universe))
	for i, j := range universe {
		stock[i] = readCSVFile(strings.Join([]string{j, "csv"}, "."))
	}
	return stock
}*/
