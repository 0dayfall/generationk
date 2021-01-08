package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// Define a type named "assets" as a slice of string
type universe []string

// Now, for our new type, implement the two methods of
// the flag.Value interface...
// The first method is String() string
func (a *universe) String() string {
	return fmt.Sprintf("%s", *a)
}

// The second method is Set(assetName string) error
func (a *universe) Set(assetName string) error {
	fmt.Printf("%s\n", assetName)
	*a = append(*a, assetName)
	return nil
}

var universeFlag universe

func init() {
	flag.Var(&universeFlag, "Universe", "Symbols to process e.g. AAPL, TSLA, V, X. Should be a comma separated list. Can also be ALL for all files in a directory.")
}

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s[post_count]\n", os.Args[0])
	flag.PrintDefaults()
}

/*
genk := NewGenerationK()
			genk.SetPortfolio(portfolio)
			genk.AddStrategy(strategy)

			now := time.Now()
			start := now.AddDate(-15, -9, -2)
			genk.SetStartDate(start)
			now = time.Now()
			end := now.AddDate(0, -3, -2)
			genk.SetEndDate(end)

			//genk.RunEventBased()
			dataManager := NewCSVDataManager(genk)
			//dataManager.SetHandler(genk)
			//genk.AddDataManager(dataManager)

			//dataManager.ReadCSVFilesAsync([]string{"test/data/ABB.csv", "test/data/ASSAb.csv"})
			ohlc := dataManager.ReadCSVFile(localFilename)
*/

func main() {
	var folderFlag string
	//var fileFormatFlag string
	var startDate string
	var endDate string
	var parallel bool

	flag.StringVar(&folderFlag, "Folder", ".", "Directory with data files.")
	//flag.FileFormatVar(&folderFlag, "Folder", ".", "Directory with data files.")
	flag.StringVar(&startDate, "StartDate", "2001-01-02 15:04:05 -0700", "Date and time of the start of the backtest.")
	flag.StringVar(&endDate, "EndDate", "2020-01-02 15:04:05 -0700", "Date and time of the end of the backtest.")
	flag.BoolVar(&parallel, "Parallel", true, "If set to false, files in the universe will be processed not in parallel as separate input, but as in one batch.")

	flag.Parse()

	layout := "2006-01-02 15:04:05 -0700"
	t1, err := time.Parse(layout, startDate)
	if err != nil {
		fmt.Println("Error while parsing date :", err)
		os.Exit(1)
	}
	fmt.Print(t1)

	t2, err := time.Parse(layout, endDate)
	if err != nil {
		fmt.Println("Error while parsing date :", err)
		os.Exit(1)
	}
	fmt.Print(t2)

	if len(flag.Args()) == 0 {
		usage()
		os.Exit(0)
	}

	//	genk.GetAssets()
}
