// +build test
package generationk

import (
	genk "generationk/internal"
	"testing"
)

func TestDatamanager(t *testing.T) {
	ctx := genk.NewContext()
	dm := genk.NewCSVDataManager(ctx)
	asset := dm.ReadCSVFile("ABB.csv")
	v := len(asset.CloseArray())
	want := 9776

	if got := v; got != want {
		t.Errorf("ReadCSVFile(\"ABB.csv\") = %d, want %d", got, want)
	}
}
