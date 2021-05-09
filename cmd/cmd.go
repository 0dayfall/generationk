package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	K "github.com/0dayfall/generationk"
	D "github.com/0dayfall/generationk/data"
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
	backtestHeaders := backtestCmd.Bool("headers", false, "CSV has headers")
	backtestReverse := backtestCmd.Bool("reverse", false, "Read file reverse")
	backtestDir := backtestCmd.String("dir", "", "Directory name")
	backtestFromDate := backtestCmd.String("fromDate", "01/01/2015", "From date")
	backtestToDate := backtestCmd.String("toDate", time.Now().Format("02/01/2006"), "To date")
	backtestCash := backtestCmd.String("cash", "0", "How much money is there")

	help := "Add the subcommand 'hurst' or 'backtest'"

	// Switch on subcommands, then apply the desired set of flags.
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, help)
		os.Exit(0)
	}

	switch os.Args[1] {
	case hurstCmd.Name():
		err := hurstCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

	case backtestCmd.Name():
		err := backtestCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

		ctx := K.NewContext()

		startDate, err := time.Parse("1/2/2006", *backtestFromDate)
		if err != nil {
			log.Fatal("Could not parse date format in ", *backtestFromDate)
		}
		ctx.SetStartDate(startDate)

		//Its OK to not set the to date
		endDate, err := time.Parse("02/01/2006", *backtestToDate)
		if err != nil {
			log.Fatal("Could not parse date format in ", *backtestToDate)
		}
		ctx.SetEndDate(endDate)

		fmt.Println(*backtestDir)
		ctx.SetDataPath(*backtestDir)

		var dm *D.DataManager

		switch *backtestMapping {

		case "investing":
			fmt.Println("Using investing.com mapping")
			dm = D.NewCSVDataManager(*backtestHeaders, *backtestReverse, D.MapRecordsInvesting)

		case "spotlight":
			fmt.Println("Using spotlight..se mapping")
			dm = D.NewCSVDataManager(*backtestHeaders, *backtestReverse, D.MapRecordsSpotlight)

		default:
			fmt.Println("Using default mapping")
			dm = D.NewCSVDataManager(*backtestHeaders, *backtestReverse, nil)

		}

		switch *backtestFile {

		case "MACrossStrategy":
			ctx.SetStrategy(new(S.MACrossStrategy))
			K.RunPlain(ctx, dm)

		case "DonchStrategy":
			ctx.SetStrategy(new(S.DonchStrategy))

			cash, err := strconv.ParseFloat(*backtestCash, 32)
			if err != nil {
				log.Fatal("Cash doesnt seem to be a number")
			}

			K.RunParallell(ctx, dm, cash)

		case "RMICrossStrategy":
			ctx.SetStrategy(new(S.RMICross))
			K.RunPlain(ctx, dm)

		case "RMIHigh":
			ctx.SetStrategy(new(S.RMIHigh))
			K.RunPlain(ctx, dm)

		case "MonthlyStrat":
			ctx.SetStrategy(new(S.MonthlyStrat))

			cash, err := strconv.ParseFloat(*backtestCash, 32)
			if err != nil {
				log.Fatal("Cash doesnt seem to be a number")
			}

			K.RunParallell(ctx, dm, cash)

		case "RebalanceStrategy":
			fmt.Printf("Rebalancing running with parameters %s, %t, %t, %s\n", *backtestMapping, *backtestHeaders, *backtestReverse, *backtestDir)
			ctx.SetStrategy(new(S.RebalanceStrat))

			cash, err := strconv.ParseFloat(*backtestCash, 32)
			if err != nil {
				log.Fatal("Cash doesnt seem to be a number")
			}

			K.RunParallell(ctx, dm, cash)

		default:
			log.Fatal("Could not find a strategy with that name in /strategies")

		}

		fmt.Println("Running job")

	default:
		log.Fatal("Example usage: backtest -test RMIStrategy -dir ..\\test\\data\\CSV2 -fromDate 01/01/2015")

	}
}
