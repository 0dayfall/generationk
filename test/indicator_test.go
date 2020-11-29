package generationk

import (
	"fmt"
	"generationk/indicators"
	genk "generationk/internal"
	"testing"
)

func TestInidcator(t *testing.T) {
	fmt.Printf("Running GENERATIONK")
	dm := genk.NewDataManager()

	asset := dm.ReadCSVFile("ABB.csv")
	v := len(asset.CloseArray())

	ma50 := indicators.SimpleMovingAverage(asset.CloseArray(), 50)
	ma200 := indicators.SimpleMovingAverage(asset.CloseArray(), 200)
	//ma = &indicators.Average{&indicators.IndicatorStruct{}}

	//&Job{command, log.New(os.Stderr, "Job: ", log.Ldate)}
	//ma50 := ma.Sma(asset.GetCloseArray(), 50)
	//ma200 := ma.Sma(asset.GetCloseArray(), 200)
	//fmt.Println("%v", ma200)

	largerThan := ma50.LargerThanInd(ma200)

	fmt.Printf("Length of bool IndicatorStruct: %v", len(largerThan))
	for _, j := range largerThan {
		if j {
			//fmt.Printf("MA50 > MA200\n")
		}

	}

	fmt.Printf("Length asset %d", v)
}
