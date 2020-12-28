package generationk_test

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	K "github.com/0dayfall/generationk"
	indicators "github.com/0dayfall/generationk/indicators"

	log "github.com/sirupsen/logrus"
)

const (
	PERCENT_LIMIT = 0.03
)

//Strategy strategy
type MA2_MA5 struct {
	ma5 indicators.SimpleMovingAverage
	ma2 indicators.SimpleMovingAverage
}

//Setup is used to declare what indicators will be used
func (ma *MA2_MA5) Once(ctx *K.Context) error {
	//ma2 is used for moving average 2
	ma.ma2 = indicators.NewSimpleMovingAverage(indicators.Close, 2)
	//ma5 is used for moving average 5
	ma.ma5 = indicators.NewSimpleMovingAverage(indicators.Close, 5)

	//Add indicators to context
	ctx.AddIndicator(&ma.ma2)
	ctx.AddIndicator(&ma.ma5)

	//The minimum data needed to calculate the moving average
	ctx.SetInitPeriod(5)

	return nil
}

//Update gets called when updates arrive
func (ma *MA2_MA5) Update(ctx *K.Context) {
	ctx.K++
}

//Tick get called when there is new data coming in
func (ma *MA2_MA5) PerBar(genkC K.GenkCallback) {
	ma2 := ma.ma2.ValueAtIndex(0)
	ma5 := ma.ma5.ValueAtIndex(0)
	diff := ma5 / ma2
	//If the diff is larger than the PERCENT_LIMIT then try to buy
	if diff-1 > PERCENT_LIMIT {
		//If we dont own the asset, then we can process the logic
		if !genkC.IsOwning(genkC.Assets()[0]) {
			//Send an order to buy asset
			err := genkC.OrderSend(genkC.Assets()[0], K.BuyOrder, K.MarketOrder, 0, 100)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if ma.ma2.ValueAtIndex(0) > ma.ma5.ValueAtIndex(0) {
		if genkC.IsOwning(genkC.Assets()[0]) {
			err := genkC.OrderSend(genkC.Assets()[0], K.SellOrder, K.MarketOrder, 0, 100)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

//OrderEvent gets called on order events
func (ma *MA2_MA5) OrderEvent(orderEvent K.Event) {
	log.WithFields(log.Fields{
		"orderEvent": orderEvent,
	}).Info("MAStrategy_test> OrderEvent")
}

func readMAFolder(folderPath string) {
	files, err := filepath.Glob(folderPath + "*.csv")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("files %s", files)
	//d.ReadCSVFilesAsync(files)
	portfolio := K.NewPortfolio()
	portfolio.SetBalance(100000)

	var wg sync.WaitGroup

	y := 0

	for _, fileName := range files {
		wg.Add(1)
		go func(localFilename string) {
			genk := K.NewGenerationK()
			genk.AddPortfolio(portfolio)

			strategy := new(MA2_MA5)

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
			dataManager := K.NewCSVDataManager(genk)
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

func TestMA2MA5Run(t *testing.T) {
	//t.Parallel()
	readMAFolder("test/data/")
}
