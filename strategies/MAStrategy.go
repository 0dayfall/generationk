package strategies

import (
	"log"

	K "github.com/0dayfall/generationk"
	I "github.com/0dayfall/generationk/indicators"
)

//Strategy strategy
type MACrossStrategy struct {
	ma50  []float64
	close []float64
}

func (ma *MACrossStrategy) GetParams() []*K.Params {
	return nil
}

//Setup is used to declare what indicators will be used
func (ma *MACrossStrategy) Once(ctx *K.Context, ohlc *K.OHLC) error {

	//The closing prices
	ma.close = ohlc.Close

	//The Simple Moving Average length 50 periods, the ones from 0 to 50 will be registred in the array as well
	ma.ma50 = I.SimpleMovingAverage(ohlc.Close, 50)

	//If the init period is set PerBar will not be called until the InitPeriod is reached
	ctx.SetInitPeriod(50)

	return nil
}

//Tick get called when there is new data coming in
func (ma *MACrossStrategy) PerBar(k int, callback K.Callback) error {

	/*fmt.Printf("k: %d ", k)
	fmt.Printf("len(close): %d ", len(ma.close))
	fmt.Printf("len(ma5) %d\n", len(ma.ma50))*/

	if ma.close[k] > ma.ma50[k] {
		if !callback.Owning() {
			err := callback.SendOrder(K.BuyOrder, K.MarketOrder, 100)

			if err != nil {
				log.Fatal(err)

				return err
			}
		}
	}

	if ma.close[k] < ma.ma50[k] {
		if callback.Owning() {
			err := callback.SendOrder(K.SellOrder, K.MarketOrder, 100)

			if err != nil {
				log.Fatal(err)

				return err
			}
		}
	}

	return nil
}

//Update is called before perBar
func (ma *MACrossStrategy) Update(k *int) error { return nil }

//OrderEvent gets called on order events
func (ma *MACrossStrategy) OrderEvent(orderEvent K.Event) {}

/*func readFolder(folderPath string) {
	files, err := filepath.Glob(folderPath + "*.csv")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("files %s", files)
	//d.ReadCSVFilesAsync(files)
	portfolio := K.NewPortfolio()
	portfolio.SetBalance(100000)

	var wg sync.WaitGroup

	y := 0

	for _, fileName := range files {
		wg.Add(1)
		go func(localFilename string, wg *sync.WaitGroup) {
			genk := K.NewGenerationK()
			genk.SetPortfolio(portfolio)
			genk.AddStrategy(new(MACrossStrategy))

			now := time.Now()
			start := now.AddDate(-15, -9, -2)
			genk.SetStartDate(start)
			now = time.Now()
			end := now.AddDate(0, -3, -2)
			genk.SetEndDate(end)

			//genk.RunEventBased()
			//dataManager := K.NewCSVDataManager()
			//dataManager.SetHandler(genk)
			//genk.AddDataManager(dataManager)

			//dataManager.ReadCSVFilesAsync([]string{"test/data/ABB.csv", "test/data/ASSAb.csv"})
			asset := K.ReadCSVFile(localFilename, false, nil)

			fmt.Printf("Local file name: %s\n\n", localFilename)
			genk.AddAsset(asset)
			runErr := genk.Run()
			if runErr != nil {
				fmt.Print(runErr.Error())
				os.Exit(0)
			}

			wg.Done()
		}(fileName, &wg)
		y++
	}
	wg.Wait()

	fmt.Printf("balance %f\n", portfolio.GetBalance())
}*/

/*func TestRun(t *testing.T) {
	//defer profile.Start().Stop()
	//t.Parallel()
	readFolder("../data/CSV2/")
}*/
