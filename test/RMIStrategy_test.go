package generationk_test

import (
	"testing"

	K "github.com/0dayfall/generationk"
	indicators "github.com/0dayfall/generationk/indicators"
	"github.com/rs/zerolog/log"
)

//Strategy strategy
type RMICrossStrategy struct {
	fastRMI   []float64
	slowRMI   []float64
	buyParam  K.Params
	sellParam K.Params
}

/*func (rmi *RMICrossStrategy) GetParams() []*K.Params {
	return []*K.Params{&rmi.buyParam}
}*/

func (rmi *RMICrossStrategy) GetParams() []*K.Params {
	//Test which of the buy params are the best ones
	rmi.buyParam = K.Params{
		Low:  -30,
		High: -20,
	}

	//Test in combination with the selling params
	rmi.sellParam = K.Params{
		Low:  10,
		High: 30,
	}

	return []*K.Params{&rmi.buyParam, &rmi.sellParam}
}

//Setup is used to declare what indicators will be used
func (rmi *RMICrossStrategy) Once(ctx *K.Context, ohlc *K.OHLC) error {

	//The Simple Moving Average length 50 periods, the ones from 0 to 50 will be registred in the array as well
	rmi.fastRMI, rmi.slowRMI = indicators.RMI(ohlc.Close, 30, 2, 250, 20)

	//If the init period is set PerBar will not be called until the InitPeriod is reached
	ctx.SetInitPeriod(50)

	return nil
}

//Tick get called when there is new data coming in
func (rmi *RMICrossStrategy) PerBar(k int, callback K.Callback) error {

	//fmt.Printf("k: %d ", k)
	//fmt.Printf("len(close): %d ", len(rmi.close))
	//fmt.Printf("len(rmi5) %d\n", len(rmi.rmi50))

	/*fmt.Printf("Buy %f\n", rmi.buyParam.GetValue())
	fmt.Printf("Sell %f\n", rmi.sellParam.GetValue())*/

	if rmi.fastRMI[k] < rmi.buyParam.GetValue() {
		if !callback.Owning() {
			err := callback.SendOrder(K.BuyOrder, K.MarketOrder, 100)

			if err != nil {
				log.Fatal().
					Err(err)

				return err
			}
		}
	}

	if rmi.fastRMI[k] > rmi.sellParam.GetValue() {
		if callback.Owning() {
			err := callback.SendOrder(K.SellOrder, K.MarketOrder, 100)

			if err != nil {
				log.Fatal().
					Err(err)

				return err
			}
		}
	}

	return nil
}

//Update is called before perBar
func (rmi *RMICrossStrategy) Update(k *int) error { return nil }

//OrderEvent gets called on order events
func (rmi *RMICrossStrategy) OrderEvent(orderEvent K.Event) {}

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

func TestRMI(t *testing.T) {
	//var cross RMICrossStrategy
	K.Run("data/", new(RMICrossStrategy))
}
