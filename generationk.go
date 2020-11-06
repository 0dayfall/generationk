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
  "sync"
)

// Stock data type
type Stock struct {
	name string
	ohlc []OHLC
}

type Strategy struct {}

func (m Strategy) init() {
  fmt.Printf("Init market\n")
}

func (m Strategy) tick(ohlc OHLC) {
  fmt.Printf("Closen %f", ohlc.close)
}

//Market for this backtester
type Market struct {
  strategy Strategy
  datamanager DataManager
}

func (m Market) setStrategy(strategy *Strategy) {
  m.strategy = *strategy
}

func (m Market) setDatamanager(datamanager *DataManager) {
  fmt.Printf("Datamanager: %d", len(datamanager.stock.ohlc))
  m.datamanager = *datamanager
  fmt.Printf("Datamanager set, len %d\n",len(m.datamanager.stock.ohlc))
}

//Run backtester
func (m Market) Run(wg *sync.WaitGroup) {
  defer wg.Done()
  fmt.Printf("Defer done\n")
  m.strategy.init()
  fmt.Printf("Len: %d\n",len(m.datamanager.stock.ohlc))
	for i,f := range m.datamanager.stock.ohlc {
    fmt.Printf("tick: %d\n", i)
    m.strategy.tick(f)
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
  stock Stock
}

func (d DataManager) readCSVFile(file string) {
	d.stock.name = filepath.Base(file)
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
	d.stock.ohlc = s
  fmt.Printf("%d\n", len(d.stock.ohlc))
}

func main() {
  var wg sync.WaitGroup
  //folderPtr := flag.String("folder", "", "a path pointing to a folder with cvs files")
	filePtr := flag.String("file", "", "a path pointing to a cvs file")
	//uniPtr := flag.String("universe", "", "abbreviation of a stock")
	flag.Parse()
	start := time.Now()
	//universe := []string{*uniPtr}
	//y := readCSV(universe)
  var dataManager DataManager
  var strategy Strategy
  var market Market

  dataManager.readCSVFile(*filePtr)
  fmt.Printf("len %d", len(dataManager.stock.ohlc))
  market.setDatamanager(&dataManager)
  market.setStrategy(&strategy)
  wg.Add(1)
  go market.Run(&wg)
  wg.Wait()
  fmt.Println(time.Since(start))
}
