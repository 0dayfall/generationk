package generationk

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Maprecords func(...string) (time.Time, []float64)

type DataManager struct {
	Folder      string
	Headers     bool
	MappingFunc Maprecords
}

func NewCSVDataManager(folder string, headers bool, mapping Maprecords) *DataManager {
	return &DataManager{
		Folder:      folder,
		Headers:     headers,
		MappingFunc: mapping,
	}
}

func (d *DataManager) ReadCSVFile(file string) *Asset {

	csvfile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}

	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()

	if err != nil && errors.Is(err, io.EOF) {
		log.Fatal(err)
	}

	asset := d.createAsset(file, records)

	return asset
}

//pasetFloat is used to parse the floats from the CSV files and is a better way to
//to handle errors
func ParseFloat(value string) float64 {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal(err)
	}

	return floatValue
}

func (d *DataManager) createAsset(file string, records [][]string) *Asset {

	size := len(records)

	var ohlc OHLC
	ohlc.Time = make([]time.Time, size)
	ohlc.Open = make([]float64, size)
	ohlc.High = make([]float64, size)
	ohlc.Low = make([]float64, size)
	ohlc.Close = make([]float64, size)
	ohlc.Volume = make([]float64, size)

	if d.MappingFunc == nil {
		d.MappingFunc = mapRecords //Use the default mapping
	}

	h := 0
	//Headers are used; dont read the last line = 0
	if d.Headers {
		h = 1
	}
	for i := size - 1; i >= h; i-- {

		time, record := d.MappingFunc(records[i]...)

		ohlc.Time[i] = time
		ohlc.Open[i] = record[0]
		ohlc.High[i] = record[1]
		ohlc.Low[i] = record[2]
		ohlc.Close[i] = record[3]
		ohlc.Volume[i] = record[4]
	}

	assetName := strings.TrimSuffix(filepath.Base(file), path.Ext(file))

	return NewAsset(assetName, &ohlc, size)
}

func mapRecords(records ...string) (time.Time, []float64) {
	time, err := time.Parse("1/2/2006 15:04:05", records[0]+" "+records[1])
	//record1, err := time.Parse("2006-01-02 15:04:05", records[i][0])
	if err != nil {
		log.Fatal(err)
	}

	floats := make([]float64, 5)

	floats[0] = ParseFloat(records[2])
	floats[1] = ParseFloat(records[3])
	floats[2] = ParseFloat(records[4])
	floats[3] = ParseFloat(records[5])
	floats[4] = ParseFloat(records[6])

	return time, floats
}

//ReadCSVFile reads a CSV file and maps the records according to this method
/*func ReadCSVFile(file string) *Asset {

	csvfile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}

	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()

	if err != nil && errors.Is(err, io.EOF) {
		log.Fatal(err)
	}

	//return mapping(file, records), nil
	size := len(records)
	var ohlc OHLC
	ohlc.Time = make([]time.Time, size)
	ohlc.Open = make([]float64, size)
	ohlc.High = make([]float64, size)
	ohlc.Low = make([]float64, size)
	ohlc.Close = make([]float64, size)
	ohlc.Volume = make([]float64, size)

	//Columns are used
	for i := size - 1; i >= 1; i-- {
		// Read each record from csv
		//record1, err := time.Parse("1/2/2006 00:00:00", records[i][0]+" "+records[i][1])
		record1, err := time.Parse("2006-01-02 15:04:05", records[i][0])
		if err != nil {
			log.Fatal(err)
		}

		record2 := parseFloat(records[i][1])
		record3 := parseFloat(records[i][2])
		record4 := parseFloat(records[i][3])
		record5 := parseFloat(records[i][4])
		record6 := parseFloat(records[i][5])

		ohlc.Time[i] = record1
		ohlc.Open[i] = record2
		ohlc.High[i] = record3
		ohlc.Low[i] = record4
		ohlc.Close[i] = record5
		ohlc.Volume[i] = record6
	}

	assetName := strings.TrimSuffix(filepath.Base(file), path.Ext(file))

	return NewAsset(assetName, &ohlc, size)
}*/
