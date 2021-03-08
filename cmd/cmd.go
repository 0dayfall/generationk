package main

import (
	"flag"
	"generationk/strategies"
	"log"
	"os"
	"time"

	K "github.com/0dayfall/generationk"
)

func main() {

	hurstCmd := flag.NewFlagSet("hurst", flag.ExitOnError)
	/*hurstDir := hurstCmd.String("dir", "", "Directory")
	hurstFromDate := hurstCmd.String("fromDate", "01/01/2015", "From date")
	hurstToDate := hurstCmd.String("toDate", time.Now().Format("02/01/2006"), "To date")*/

	backtestCmd := flag.NewFlagSet("backtest", flag.ExitOnError)
	backtestFile := backtestCmd.String("test", "", "Name of the struct with backtest")
	backtestDir := backtestCmd.String("dir", "", "Directory name")
	backtestFromDate := backtestCmd.String("fromDate", "01/01/2015", "From date")
	backtestToDate := backtestCmd.String("toDate", time.Now().Format("02/01/2006"), "To date")

	switch os.Args[1] {
	case hurstCmd.Name():
		hurstCmd.Parse(os.Args[2:])

		os.Exit(0)

	case backtestCmd.Name():
		backtestCmd.Parse(os.Args[2:])

		ctx := K.NewContext()

		startDate, err := time.Parse("1/2/2006", *backtestFromDate)
		if err != nil {
			log.Fatal("Could not parse date format in ", *backtestFromDate)
		}
		ctx.SetStartDate(startDate)

		//Its OK to not set the to date
		endDate, err := time.Parse("1/2/2006", *backtestToDate)
		if err != nil {
			log.Fatal("Could not parse date format in ", *backtestToDate)
		}
		ctx.SetEndDate(endDate)

		ctx.SetDataPath(*backtestDir)

		switch *backtestFile {

		case "MACrossStrategy":
			break

		case "RMIStrategy":
			dataManager := K.DataManager{
				Folder: *backtestDir,
				//Folder:      "../data/CSV1/",
				MappingFunc: nil,
			}
			K.Run(ctx, dataManager, new(strategies.RMICrossStrategy))

			break

		default:
			log.Fatal("Could not find a strategy with that name in /strategies")
		}

	default:
		log.Fatal("Please use any of the commands")

	}
}
