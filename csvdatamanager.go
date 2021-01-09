package generationk

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
		fmt.Println(err)
	}

	return floatValue
}

//NewCSVDataManager creates a new data manager object
func NewCSVDataManager(dataHandler DataHandler) *CSVDataManager {

	dm := &CSVDataManager{
		callback: dataHandler,
	}

	return dm
}

func (d *CSVDataManager) SetHandler(dataHandler DataHandler) {
	d.callback = dataHandler
}

//ReadCSVFile reads a CSV file and maps the records according to this method
func (d *CSVDataManager) ReadCSVFile(file string) (*Asset, error) {

	csvfile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}

	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()

	if err != nil && errors.Is(err, io.EOF) {
		fmt.Println(err)
		return nil, err
	}

	size := len(records)
	fmt.Printf("size: %d\n\n", size)
	var ohlc OHLC
	ohlc.Time = make([]time.Time, size)
	ohlc.Open = make([]float64, size)
	ohlc.High = make([]float64, size)
	ohlc.Low = make([]float64, size)
	ohlc.Close = make([]float64, size)
	ohlc.Volume = make([]float64, size)

	for i := size - 1; i >= 0; i-- {
		// Read each record from csv
		record1, err := time.Parse("1/2/2006 00:00:00", records[i][0]+" "+records[i][1])
		if err != nil {
			fmt.Printf("\n\n%s: %v\n", file, err)
			return nil, err
		}

		record2 := parseFloat(records[i][2])
		record3 := parseFloat(records[i][3])
		record4 := parseFloat(records[i][4])
		record5 := parseFloat(records[i][5])
		record6 := parseFloat(records[i][6])

		ohlc.Time[i] = record1
		ohlc.Open[i] = record2
		ohlc.High[i] = record3
		ohlc.Low[i] = record4
		ohlc.Close[i] = record5
		ohlc.Volume[i] = record6
	}
	assetName := strings.TrimSuffix(filepath.Base(file), path.Ext(file))

	return NewAsset(assetName, &ohlc, size), nil
}
