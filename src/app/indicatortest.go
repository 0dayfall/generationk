package main

import (
	"fmt"
	"indicators"
)

func main() {
	var dm DataManager

	asset := dm.ReadCSVFile("../../ABB.csv")

	var ma *indicators.Indicator
	ma = &indicators.Average{&indicators.Comparison{}}

	//&Job{command, log.New(os.Stderr, "Job: ", log.Ldate)}
	ma50 := ma.Sma(asset.GetCloseArray(), 50)
	ma200 := ma.Sma(asset.GetCloseArray(), 200)
	largerThan := ma.LargerThanInd(ma)

	for i, j := range largerThan {
		if j[i] {
			fmt.Printf("MA50 > MA200\n")
		} else {
			fmt.Printf("Ej större\n")
		}
	}
}
