package generationk

import (
	"testing"
	genk "generationk/internal"
	"time"
)

func TestContext(t *testing.T) {
	var c genk.Context
	dm := genk.NewDataManager()
	
	abb := dm.ReadCSVFile("ABB.csv")
	c.AddAsset(&abb)
	c.AddStrategy(&genk.Strategy{})
	c.AddStartDate(time.Now())
	c.AddEndDate(time.Now())

	v := len(c.Asset)
	want := 1
	
	if got := v; got != want {
		t.Errorf("ReadCSVFile(\"ABB.csv\") = %d, want %d", got, want)
	}
}