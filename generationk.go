package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Stock data type
type Stock struct {
	name string
	ohlc []OHLC
}

//Market for this backtester
type Market struct {
	stocks []Stock
}

//Run backtester
func (s Market) Run() {

	for i := range s.stocks {
		for j := range s.stocks[i]
	}
}

// OHLC data type
type OHLC struct {
	time                   time.Time
	open, high, low, close float64
	volume                 int
}

// Portfolio structure
type Portfolio struct {
}

//Event type
type Event struct{}

//Tick event type
type Tick Event

//Signal event type
type Signal Event

//Order event type
type Order Event

//Fill event type
type Fill Event

/*func readCSV(universe []string) []Stock {
	stock := make([]Stock, len(universe))
	for i, j := range universe {
		stock[i] = readCSVFile(strings.Join([]string{j, "csv"}, "."))
	}
	return stock
}*/

type DataManager struct {

}

func (d DataManager) readCSVFile(file string) Stock {
	var stock Stock
	stock.name = filepath.Base(file)
	csvfile, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}

	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)

	// Iterate through the records
	i := 0
	s := make([]OHLC, 1000)
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		record1, err := time.Parse("1/2/2006 00:00:00", record[0]+" "+record[1])
		record2, err := strconv.ParseFloat(record[2], 64)
		record3, err := strconv.ParseFloat(record[3], 64)
		record4, err := strconv.ParseFloat(record[4], 64)
		record5, err := strconv.ParseFloat(record[5], 64)
		record6, err := strconv.Atoi(record[6])

		if err != nil {
			log.Fatal(err)
		}

		ohlc := OHLC{time: record1, open: record2, high: record3, low: record4, close: record5, volume: record6}

		if i < len(s) {
			s[i] = ohlc
		} else {
			s = append(s, ohlc)
		}
		i++
		//fmt.Printf("In addValue: s is %v\n", s)
	}
	stock.ohlc = s
	return stock
}

func main() {
	eventQueue := make(chan Event)
	//folderPtr := flag.String("folder", "", "a path pointing to a folder with cvs files")
	filePtr := flag.String("file", "", "a path pointing to a cvs file")
	//uniPtr := flag.String("universe", "", "abbreviation of a stock")
	flag.Parse()
	start := time.Now()
	//universe := []string{*uniPtr}
	//y := readCSV(universe)
	s := readCSVFile(*filePtr)
	fmt.Println(time.Since(start))
	fmt.Printf("%v", s)
	portfolio(eventQueue)
}

func portfolio(eventQueue chan Event) {
	event := <-eventQueue
}
