package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	K "github.com/0dayfall/generationk"
	S "github.com/0dayfall/generationk/strategies"
)

func main() {

	hurstCmd := flag.NewFlagSet("hurst", flag.ExitOnError)
	/*hurstDir := hurstCmd.String("dir", "", "Directory")
	hurstFromDate := hurstCmd.String("fromDate", "01/01/2015", "From date")
	hurstToDate := hurstCmd.String("toDate", time.Now().Format("02/01/2006"), "To date")*/

	backtestCmd := flag.NewFlagSet("backtest", flag.ExitOnError)
	backtestFile := backtestCmd.String("test", "", "Name of the struct with backtest")
	backtestMapping := backtestCmd.String("mapping", "", "Mapping function")
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

		var dm *K.DataManager

		switch *backtestMapping {

		case "investing":
			fmt.Println("Using investing.com mapping")
			dm = K.NewCSVDataManager(*backtestDir, true, K.MapRecordsInvesting)

		default:
			fmt.Println("Using default mapping")
			dm = K.NewCSVDataManager(*backtestDir, false, nil)

		}

		switch *backtestFile {

		case "MACrossStrategy":
			ctx.SetStrategy(new(S.MACrossStrategy))
			K.Run(ctx, dm)

		case "RMICrossStrategy":
			ctx.SetStrategy(new(S.RMICrossStrategy))
			K.Run(ctx, dm)

		case "RebalanceStrategy":
			ctx.SetStrategy(new(S.RebalanceStrat))
			K.Run(ctx, dm)

		default:
			log.Fatal("Could not find a strategy with that name in /strategies")

		}

	default:
		log.Fatal("Example usage: backtest -test RMIStrategy -dir ..\\test\\data\\CSV2 -fromDate 01/01/2015")

	}
}
