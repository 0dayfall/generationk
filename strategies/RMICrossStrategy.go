package generationk

import (
	"log"

	K "github.com/0dayfall/generationk"
	D "github.com/0dayfall/generationk/data"
	I "github.com/0dayfall/generationk/indicators"
)

//Strategy strategy
type RMICross struct {
	fastRMI   []float64
	slowRMI   []float64
	buyParam  K.Params
	sellParam K.Params
}

func (donch *RMICross) GetParams() []*K.Params { return nil }

//Setup is used to declare what indicators will be used
func (rmi *RMICross) Once(ctx *K.Context, asset []*D.Asset) error {
	ohlc := asset[0].Ohlc

	//The Simple Moving Average length 50 periods, the ones from 0 to 50 will be registred in the array as well
	rmi.fastRMI, rmi.slowRMI = I.RMI(ohlc.Close, 30, 2, 250, 20)

	//If the init period is set PerBar will not be called until the InitPeriod is reached
	ctx.SetInitPeriod(50)

	return nil
}

//Tick get called when there is new data coming in
func (rmi *RMICross) PerBar(k int, callback K.Callback) error {

	//fmt.Printf("k: %d ", k)
	//fmt.Printf("len(close): %d ", len(rmi.close))
	//fmt.Printf("len(rmi5) %d\n", len(rmi.rmi50))

	/*fmt.Printf("Buy %f\n", rmi.buyParam.GetValue())
	fmt.Printf("Sell %f\n", rmi.sellParam.GetValue())*/

	if rmi.fastRMI[k] < rmi.buyParam.GetValue() {
		if !callback.Owning() {
			_, err := callback.SendOrder(K.BuyOrder, K.MarketOrder, 100)

			if err != nil {
				log.Fatal(err)

				return err
			}
		}
	}

	if rmi.fastRMI[k] > rmi.sellParam.GetValue() {
		if callback.Owning() {
			_, err := callback.SendOrder(K.SellOrder, K.MarketOrder, 100)

			if err != nil {
				log.Fatal(err)

				return err
			}
		}
	}

	return nil
}

func (rmi *RMICross) End(k int, callback K.Callback) error { return nil }

//Update is called before perBar
func (rmi *RMICross) Update(k *int) error { return nil }

//OrderEvent gets called on order events
func (rmi *RMICross) OrderEvent(orderEvent K.Event) {}

/*func readOneFolder(folderPath string) {
	files, err := filepath.Glob(folderPath + "*.csv")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("files %s", files)

	for _, fileName := range files {
		wg.Add(1)

		go run(fileName, &wg, c)
		y++
	}

	//Close channel when go threads are done
	go func(ch chan float64, wg *sync.WaitGroup) {
		wg.Wait()
		close(c)
	}(c, &wg)

	//Hang here until results are ready = channel closed
	for {
		fmt.Println(<-c)
	}
}*/

/*func TestRMI(t *testing.T) {
	//var cross RMICrossStrategy
	dataManager := K.DataManager{
		Folder:      "../data/CSV1/",
		MappingFunc: mapping,
	}
	K.Run(dataManager, new(RMICrossStrategy))

}

func mapping(file string, headers bool, records [][]string) *K.Asset {
	//return mapping(file, records), nil
	size := len(records)
	var ohlc K.OHLC
	ohlc.Time = make([]time.Time, size)
	ohlc.Open = make([]float64, size)
	ohlc.High = make([]float64, size)
	ohlc.Low = make([]float64, size)
	ohlc.Close = make([]float64, size)
	ohlc.Volume = make([]float64, size)

	h := 0
	//Headers are used; dont read the last line = 0
	if headers {
		h = 1
	}
	for i := size - 1; i >= h; i-- {
		// Read each record from csv
		//record1, err := time.Parse("1/2/2006 15:04:05", records[i][0]+" "+records[i][1])
		record1, err := time.Parse("2006-01-02 15:04:05", records[i][0])
		if err != nil {
			log.Fatal(err)
		}

		record2 := ParseFloat(records[i][1])
		record3 := ParseFloat(records[i][2])
		record4 := ParseFloat(records[i][3])
		record5 := ParseFloat(records[i][4])
		record6 := ParseFloat(records[i][5])

		ohlc.Time[i] = record1
		ohlc.Open[i] = record2
		ohlc.High[i] = record3
		ohlc.Low[i] = record4
		ohlc.Close[i] = record5
		ohlc.Volume[i] = record6
	}

	assetName := strings.TrimSuffix(filepath.Base(file), path.Ext(file))

	return K.NewAsset(assetName, &ohlc, size)
}

//pasetFloat is used to parse the floats from the CSV files and is a better way to
//to handle errors
func ParseFloat(value string) float64 {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal(err)
	}

	return floatValue
}*/
