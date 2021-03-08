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

type Mapping func(file string, headers bool, records [][]string) *Asset

type DataManager struct {
	Folder      string
	MappingFunc Mapping
}

func (d *DataManager) ReadCSVFile(file string, headers bool, mapping Mapping) *Asset {

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

	if d.MappingFunc != nil {
		return d.MappingFunc(file, headers, records)
	} else {
		return d.defaultMapping(file, headers, records)
	}

}

//pasetFloat is used to parse the floats from the CSV files and is a better way to
//to handle errors
func (d *DataManager) ParseFloat(value string) float64 {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal(err)
	}

	return floatValue
}

func (d *DataManager) defaultMapping(file string, headers bool, records [][]string) *Asset {
	//return mapping(file, records), nil
	size := len(records)
	var ohlc OHLC
	ohlc.Time = make([]time.Time, size)
	ohlc.Open = make([]float64, size)
	ohlc.High = make([]float64, size)
	ohlc.Low = make([]float64, size)
	ohlc.Close = make([]float64, size)
	ohlc.Volume = make([]float64, size)

	h := 0
	//Headers are used; dont read the last line = 0
	if headers {
		h = 1
	}
	for i := size - 1; i >= h; i-- {
		// Read each record from csv
		record1, err := time.Parse("1/2/2006 15:04:05", records[i][0]+" "+records[i][1])
		//record1, err := time.Parse("2006-01-02 15:04:05", records[i][0])
		if err != nil {
			log.Fatal(err)
		}

		record2 := d.ParseFloat(records[i][2])
		record3 := d.ParseFloat(records[i][3])
		record4 := d.ParseFloat(records[i][4])
		record5 := d.ParseFloat(records[i][5])
		record6 := d.ParseFloat(records[i][6])

		ohlc.Time[i] = record1
		ohlc.Open[i] = record2
		ohlc.High[i] = record3
		ohlc.Low[i] = record4
		ohlc.Close[i] = record5
		ohlc.Volume[i] = record6
	}

	assetName := strings.TrimSuffix(filepath.Base(file), path.Ext(file))

	return NewAsset(assetName, &ohlc, size)
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
