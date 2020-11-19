package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
	genk "generationk/internal"
)

//Run backtester
func Run(m *genk.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	m.Strategy.Init(m)
	//reateGraph(m)
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
	dataManager := genk.NewDataManager()
	var strategy genk.Strategy
	var market genk.Context
	var asset genk.Asset

	asset = dataManager.ReadCSVFile(*filePtr)
	market.SetAsset(&asset)
	market.SetStrategy(&strategy)
	wg.Add(1)
	go Run(&market, &wg)
	wg.Wait()
	fmt.Println(time.Since(start))
}
