package generationk

import (
	"testing"
)

const numberOfRows = 12

func TestHeaderNewestfirst(t *testing.T) {

	//Create a data manager with the default mapping; but different format in file
	dm := NewCSVDataManager(true, true, nil)
	_, err := dm.ReadCSVFile("../test/data/test/csv_header_newestfirst.csv")
	if err == nil {
		t.Errorf(err.Error())
	}

	//Change to correct mapping function for this file format
	dm.MappingFunc = MapRecordsInvesting
	asset, err := dm.ReadCSVFile("../test/data/test/csv_header_newestfirst.csv")
	if err != nil {
		t.Errorf(err.Error())
	}

	//Check the legth, we should skip the header
	if len(asset.Ohlc.Open) != numberOfRows {
		t.Errorf("Wrong length")
	}

	//Make sure the first row is the last in the file
	if asset.Ohlc.Open[0] != 202.80 {
		t.Errorf("Read file in wrong order: open: %f is not the open of the 2nd row", asset.Ohlc.Open[0])
	}

	//Check that the last row is the 2nd row without header
	if asset.Ohlc.Close[numberOfRows-1] != 201.30 {
		t.Errorf("The last row is not the top row: close: %f was not the close of the 2nd row", asset.Ohlc.Close[numberOfRows-2])
	}
}

func TestNewestfirst(t *testing.T) {

	//Create a data manager with the default mapping; but different format in file
	//Should create an error which is correct
	dm := NewCSVDataManager(false, true, nil)
	_, err := dm.ReadCSVFile("../test/data/test/csv_newestfirst.csv")
	if err == nil {
		t.Errorf(err.Error())
	}

	//Change the mapping functino from default to the one for this type of file
	dm.MappingFunc = MapRecordsInvesting
	asset, err := dm.ReadCSVFile("../test/data/test/csv_newestfirst.csv")
	if err != nil {
		t.Errorf(err.Error())
	}

	//Check the legth, we should skip the header
	if len(asset.Ohlc.Open) != numberOfRows {
		t.Errorf("Wrong length")
	}

	//Make sure the first row is the last in the file
	if asset.Ohlc.Open[0] != 202.80 {
		t.Errorf("Read file in wrong order: open: %f is not the open of the 2nd row", asset.Ohlc.Open[0])
	}

	//Check that the last row is the 2nd row without header
	if asset.Ohlc.Close[numberOfRows-1] != 201.30 {
		t.Errorf("The last row is not the top row: close: %f was not the close of the 2nd row", asset.Ohlc.Close[numberOfRows-2])
	}
}

const numberOfRowsOldestFirst = 70

func TestOldestfirst(t *testing.T) {

	//Create a data manager with the default mapping; but different format in file
	//Should create an error which is correct
	dm := NewCSVDataManager(false, false, nil)
	asset, err := dm.ReadCSVFile("../test/data/test/csv_noheader_oldestfirst.csv")
	if err != nil {
		t.Errorf(err.Error())
	}

	//Check the legth, we should skip the header
	if len(asset.Ohlc.Open) != numberOfRowsOldestFirst {
		t.Errorf("Wrong length")
	}

	//Make sure the first row is the first in the file
	if asset.Ohlc.Open[0] != 23.12 {
		t.Errorf("Read file in wrong order: open: %f is not the open of the 1st row", asset.Ohlc.Open[0])
	}

	//Check that the last row is the last
	if asset.Ohlc.Close[numberOfRowsOldestFirst-1] != 21.5 {
		t.Errorf("The last row is not the last row: close: %f was not the close of the last row", asset.Ohlc.Close[numberOfRows-2])
	}
}

func TestPadding(t *testing.T) {

	//Create a data manager with the default mapping; but different format in file
	//Should create an error which is correct
	dm := NewCSVDataManager(false, false, nil)
	asset, err := dm.ReadCSVFile("../test/data/test/csv_noheader_oldestfirst.csv")
	if err != nil {
		t.Errorf(err.Error())
	}
	if asset.Length != 70 {
		t.Errorf("Length is not 70, length is: %d", asset.Length)
	}

	//Create a data manager with the default mapping; but different format in file
	//Should create an error which is correct
	dm2 := NewCSVDataManager(true, true, MapRecordsInvesting)
	asset2, err2 := dm2.ReadCSVFile("../test/data/test/csv_newestfirst.csv")
	if err2 != nil {
		t.Errorf(err2.Error())
	}
	if asset2.Length != 11 {
		t.Errorf("Length is not 11, length is: %d", asset2.Length)
	}

	assets := []*Asset{asset, asset2}
	dm.CreatePadding(assets)
	for i := 0; i < len(assets); i++ {
		//fmt.Printf("%v", assets[i].Ohlc.Time)
		if assets[i].Length != 70 {
			t.Errorf("assets[%d].Length = %d", i, assets[i].Length)
		}
	}

}
