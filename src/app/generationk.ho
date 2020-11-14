package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

//Run backtester
func Run(m *Context, wg *sync.WaitGroup) {
	defer wg.Done()
	m.strategy.init(m)
	createGraph(m)
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
	var market Context
	var asset Asset

	asset = dataManager.readCSVFile(*filePtr)
	market.setAsset(&asset)
	market.setStrategy(&strategy)
	wg.Add(1)
	go Run(&market, &wg)
	wg.Wait()
	fmt.Println(time.Since(start))
}
