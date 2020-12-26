package generationk

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	indicators "github.com/0dayfall/generationk/indicators"

	log "github.com/sirupsen/logrus"
)

//Strategy strategy
type MACrossStrategy struct {
	ma50  indicators.SimpleMovingAverage
	close indicators.TimeSeries
}

//Setup is used to declare what indicators will be used
func (ma *MACrossStrategy) Setup(ctx *Context) error {
	//Want access to the latest 5 closing prices
	ma.close = indicators.NewTimeSeries(indicators.Close, 5)
	//MA50
	ma.ma50 = indicators.NewSimpleMovingAverage(indicators.Close, 50)

	//Add indicators to context
	ctx.AddIndicator(&ma.close)
	ctx.AddIndicator(&ma.ma50)

	//The data needed to calculate MA
	ctx.SetInitPeriod(50)

	return nil
}

//Update gets called when updates arrive
func (ma *MACrossStrategy) Update(ctx *Context) {
	ctx.K++
}

//Tick get called when there is new data coming in
func (ma *MACrossStrategy) Tick(genkC GenkCallback) {
	/*log.WithFields(
		log.Fields{
			"Assets": genkC.Assets(),
		}).Info("Assets")

	log.WithFields(
		log.Fields{
			"MA values": ma.ma50.Values(),
		}).Info("MA")*/

	if ma.close.ValueAtIndex(0) > ma.ma50.ValueAtIndex(0) {
		if !genkC.IsOwning(genkC.Assets()[0]) {
			genkC.OrderSend(genkC.Assets()[0], OrderType(BuyOrder), 0, 100)
		}
	}

	if ma.close.ValueAtIndex(0) < ma.ma50.ValueAtIndex(0) {
		if genkC.IsOwning(genkC.Assets()[0]) {
			genkC.OrderSend(genkC.Assets()[0], OrderType(SellOrder), 0, 100)
		}
	}

}

//OrderEvent gets called on order events
func (ma *MACrossStrategy) OrderEvent(orderEvent Event) {
	log.WithFields(log.Fields{
		"orderEvent": orderEvent,
	}).Debug("MAStrategy_test> OrderEvent")
}

func readFolder(folderPath string) {

	files, err := filepath.Glob(folderPath + "*.csv")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("files %s", files)
	//d.ReadCSVFilesAsync(files)
	portfolio := NewPortfolio()
	portfolio.SetBalance(100000)

	var wg sync.WaitGroup
	y := 0
	for _, fileName := range files {
		wg.Add(1)
		go func(localFilename string) {
			genk := NewGenerationK()
			genk.AddPortfolio(portfolio)

			strategy := new(MACrossStrategy)

			//Going to run with the following data thingie to collect the data
			//assetName := strings.TrimSuffix(filepath.Base(fileName), path.Ext(fileName))
			//genk.AddAsset(NewAsset(assetName, OHLC{}))
			//genk.AddAsset(NewAsset(assetName, OHLC{}))
			genk.AddStrategy(strategy)

			//genk.SetBalance(100000)
			now := time.Now()
			start := now.AddDate(-15, -9, -2)
			genk.AddStartDate(start)
			now = time.Now()
			end := now.AddDate(0, -3, -2)
			genk.AddEndDate(end)

			//genk.RunEventBased()
			dataManager := NewCSVDataManager(genk)
			//genk.AddDataManager(dataManager)

			//dataManager.ReadCSVFilesAsync([]string{"test/data/ABB.csv", "test/data/ASSAb.csv"})
			count := dataManager.ReadCSVFile(localFilename)
			log.WithFields(log.Fields{
				"count": count,
			}).Info("Number of lines processed")

			wg.Done()
		}(fileName)
		y++
	}
	wg.Wait()
	log.WithFields(log.Fields{
		"balance": portfolio.GetBalance(),
	}).Info("Balance")
}

func TestRun(t *testing.T) {
	//files := []string{"test/data/ABB.csv", "test/data/ASSAb.csv", "test/data/BILL.csv"}

	readFolder("test/data/")
}
