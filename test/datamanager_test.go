package generationk

import (
	"testing"
	genk "generationk/internal"
)

func TestDatamanager(t *testing.T) {
	dm := genk.NewDataManager()
	asset := dm.ReadCSVFile("ABB.csv")
	v := len(asset.GetCloseArray())
	want := 9776

	if got := v; got != want {
		t.Errorf("ReadCSVFile(\"ABB.csv\") = %d, want %d", got, want)
	}
}
