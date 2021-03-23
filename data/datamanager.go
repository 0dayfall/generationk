package generationk

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Maprecords func(...string) (time.Time, []float64, error)

type DataManager struct {
	Headers     bool
	Reverse     bool
	MappingFunc Maprecords
}

func NewCSVDataManager(headers bool, reverse bool, mapping Maprecords) *DataManager {
	fmt.Println("Returning data manager")
	return &DataManager{
		Headers:     headers,
		Reverse:     reverse,
		MappingFunc: mapping,
	}
}

func (d *DataManager) ReadCSVFiles(folder string) []*Asset {
	return nil
}

func (d *DataManager) CreatePadding(asset []*Asset) {
	var length int

	//Find the length of the longest asset
	for i := 0; i < len(asset); i++ {
		if length < asset[i].Length {
			length = asset[i].Length
		}
	}

	//For each asset
	for i := 0; i < len(asset); i++ {

		//As long as this asset is shorter than the longest
		if asset[i].Length < length {

			//Store the difference between the lengths
			diff := length - asset[i].Length

			//Create empty slices with 0 values
			padTime := make([]time.Time, diff)
			padOpen := make([]float64, diff)
			padHigh := make([]float64, diff)
			padLow := make([]float64, diff)
			padClose := make([]float64, diff)
			padVolume := make([]float64, diff)

			asset[i].Ohlc.Time = append(padTime, asset[i].Ohlc.Time...)
			asset[i].Ohlc.Open = append(padOpen, asset[i].Ohlc.Open...)
			asset[i].Ohlc.High = append(padHigh, asset[i].Ohlc.High...)
			asset[i].Ohlc.Low = append(padLow, asset[i].Ohlc.Low...)
			asset[i].Ohlc.Close = append(padClose, asset[i].Ohlc.Close...)
			asset[i].Ohlc.Volume = append(padVolume, asset[i].Ohlc.Volume...)

			//Adjust the length of the asset to the new length
			asset[i].Length = len(asset[i].Ohlc.High)
		}
	}
}

//ReadCSVFile reads the file in file string and return as Asset based on the columns of
//Open, High, Low, Close, Volume
func (d *DataManager) ReadCSVFile(file string) (*Asset, error) {
	csvfile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	//Take the headers off before creating the asset
	if d.Headers {
		records = records[1:]
	}
	asset, err := d.createAsset(file, records)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (d *DataManager) createAsset(file string, records [][]string) (*Asset, error) {
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

	//We read every line
	for i := 0; i < size; i++ {

		var time time.Time
		var record []float64
		var err error
		if d.Reverse {
			//size 12: from 0 to 11
			//last row size - 0 - 1: 11
			//first row size - 10 -1: 0
			time, record, err = d.MappingFunc(records[size-i-1]...)
		} else {
			//Per usual
			time, record, err = d.MappingFunc(records[i]...)
		}

		if err != nil {
			return nil, err
		}

		ohlc.Time[i] = time
		ohlc.Open[i] = record[0]
		ohlc.High[i] = record[1]
		ohlc.Low[i] = record[2]
		ohlc.Close[i] = record[3]
		ohlc.Volume[i] = record[4]
	}

	assetName := strings.TrimSuffix(filepath.Base(file), path.Ext(file))

	return NewAsset(assetName, &ohlc, size), nil
}

func mapRecords(records ...string) (time.Time, []float64, error) {
	//Parse the time first and then each columns of the file into OHLCV
	time, err := time.Parse("1/2/2006 15:04:05", records[0]+" "+records[1])
	if err != nil {
		return time, nil, err
	}

	//Create an array of floats for O,H,L,C,V
	floats := make([]float64, 5)

	records[2] = stripThousandSeparator(records[2])
	floats[0], err = strconv.ParseFloat(records[2], 64)
	if err != nil {
		return time, nil, err
	}

	records[3] = stripThousandSeparator(records[3])
	floats[1], err = strconv.ParseFloat(records[3], 64)
	if err != nil {
		return time, nil, err
	}

	records[4] = stripThousandSeparator(records[4])
	floats[2], err = strconv.ParseFloat(records[4], 64)
	if err != nil {
		return time, nil, err
	}

	records[5] = stripThousandSeparator(records[5])
	floats[3], err = strconv.ParseFloat(records[5], 64)
	if err != nil {
		return time, nil, err
	}

	records[6] = stripThousandSeparator(records[6])
	floats[4], err = strconv.ParseFloat(records[6], 64)
	if err != nil {
		//fmt.Println(records[6])
		return time, nil, err
	}

	return time, floats, nil
}

func stripThousandSeparator(s string) string {
	return strings.Replace(s, ",", "", -1)
}
