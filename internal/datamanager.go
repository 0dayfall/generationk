package internal

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//DataManager type
type DataManager struct {
	fieldMapper func(records [][]string) []OHLC
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
func NewDataManager() DataManager {
	return newDataManager(nil)
}

//NewDataManager is used to create a data manager. It uses a field mapper to map the field
func newDataManager(mapper func(records [][]string) []OHLC) DataManager {
	var dm DataManager
	if mapper == nil {
		dm.fieldMapper = defaultFieldMapper
	} else {
		dm.fieldMapper = mapper
	}
	return dm
}

func defaultFieldMapper(records [][]string) []OHLC {
	// Iterate through the records
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

		ohlc := OHLC{time: record1, open: record2, high: record3, low: record4, close: record5, volume: record6}

		if i < len(s) {
			s[i] = ohlc
		} else {
			s = append(s, ohlc)
		}
		//fmt.Printf("In addValue: s is %v\n", s)
	}
	reverseSlice(&s)
	return s
}

func reverseSlice(ohlc *[]OHLC) {
	for i := len(*ohlc)/2 - 1; i >= 0; i-- {
		opp := len(*ohlc) - 1 - i
		(*ohlc)[i], (*ohlc)[opp] = (*ohlc)[opp], (*ohlc)[i]
	}
}

//ReadCSVFile reads a CSV file
func (d *DataManager) ReadCSVFile(file string) Asset {
	var stock Asset
	stock.name = filepath.Base(file)
	csvfile, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}

	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()

	if err != nil && err != io.EOF {
		log.Printf("Was not possible to read the file %s", err)
	}

	s := d.fieldMapper(records)
	stock.ohlc = s

	return stock
}

/*func readCSV(universe []string) []Stock {
	stock := make([]Stock, len(universe))
	for i, j := range universe {
		stock[i] = readCSVFile(strings.Join([]string{j, "csv"}, "."))
	}
	return stock
}*/
