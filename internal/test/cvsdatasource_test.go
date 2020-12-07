package internal

import (
	genk "generationk/internal"
	"testing"
)

func TestDatamanager(t *testing.T) {
	//Context that the strategy is being run with such as assets
	market := genk.NewContext()
	dm := genk.NewCSVDataManager(market)
	dm.ReadCSVFile("Assa.csv")
}
