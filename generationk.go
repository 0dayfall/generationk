package generationk

import (
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	D "github.com/0dayfall/generationk/data"
)

//var o sync.Once

var AssetDoesNotExist = errors.New("Asset does not exist")
var FFToStartDate = errors.New("Fast forwarding to start date")
var EndOfBacktest = errors.New("End of backtest")
var EndOfData = errors.New("End of data")
var UnstablePeriod = errors.New("The stable period is not yet reached")
var Initialization = errors.New("Initialization in Once() failed")
var EndOfAsset = errors.New("End of asset")

//Callback is used in the strategy to give actions back to the backtest in progress
type Callback interface {
	Owning() bool
	IsOwning(assetName string) (bool, error)
	SendOrder(direction Direction, orderType OrderType, qty int) (float64, error)
	SendOrderFor(assetName string, direction Direction, orderType OrderType, qty int) (float64, error)
	Sell(assetName string) (float64, error)
	Assets() []string
	Date() time.Time
	Record(name string, record RecordStruct)
	PositionSize(percentageOfCapital float64, price float64) int
}

type GenerationK struct {
	ctx     *Context
	records map[string][]RecordStruct
}

//NewGenerationK is used to create a new backtest
func NewGenerationK() *GenerationK {

	generationK := &GenerationK{
		ctx:     NewContext(),
		records: make(map[string][]RecordStruct),
	}

	return generationK
}

func (g *GenerationK) inc() {
	g.ctx.K++
}

type intervalFunc func(oldTime time.Time, newTime time.Time) bool

func determineInterval(interval string) intervalFunc {
	switch interval {
	case "D":

		return func(timeOld time.Time, timeNew time.Time) bool {
			_, _, dayOld := timeOld.Date()
			_, _, dayNew := timeNew.Date()

			return dayNew > dayOld
		}

	case "M":

		return func(timeOld time.Time, timeNew time.Time) bool {
			_, monthOld, _ := timeOld.Date()
			_, monthNew, _ := timeNew.Date()

			return monthNew > monthOld
		}

	case "Q":

		return func(timeOld time.Time, timeNew time.Time) bool {
			_, monthNew, _ := timeNew.Date()

			return monthNew == 3 || monthNew == 6 ||
				monthNew == 9 || monthNew == 12
		}

	}

	return nil
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (g *GenerationK) nextDate() error {
	//Increase the time with a day
	g.ctx.datePointer = g.ctx.datePointer.AddDate(0, 0, 1)

	for _, asset := range g.ctx.assets {
		if DateEqual(g.ctx.datePointer, asset.Ohlc.Time[g.ctx.K+1]) {
			return nil
		}
	}
	return nil
}

var timer intervalFunc

func (g *GenerationK) nextGen() error {
	defer g.inc()

	g.ctx.datePointer = g.ctx.asset.Ohlc.Time[g.ctx.K]

	//Have to run this first so that we dont increase k by FF
	if g.ctx.K < 1 {
		fmt.Println("Running once")
		err := g.ctx.strategy[0].Once(g.ctx, g.ctx.assets)
		if err != nil {
			return err
		}

		// Determine which function to use for rebalancing: Daily, Monthly or Quarterly
		v, ok := interface{}(g.ctx.strategy[0]).(RebalanceStrategy)
		if ok {
			timer = determineInterval(v.GetInterval())
		} else {
			timer = nil
			fmt.Println("Interface is of wrong type")
		}
	}

	if g.ctx.asset.Ohlc.Time[g.ctx.K].Before(g.ctx.startDate) {
		return FFToStartDate
	}

	if g.ctx.asset.Ohlc.Time[g.ctx.K].After(g.ctx.endDate) {
		return EndOfBacktest
	}

	//Run setup after initperiod is finished
	if g.ctx.K < g.ctx.GetInitPeriod() {
		return UnstablePeriod
	}

	// Check if the timer function is set
	// Check if the interface implementents rebalance function
	// Call the rebalance with the date as an additional parameter
	if timer != nil {
		if timer(g.ctx.asset.Ohlc.Time[g.ctx.K-1], g.ctx.asset.Ohlc.Time[g.ctx.K]) {

			v, ok := interface{}(g.ctx.strategy[0]).(RebalanceStrategy)

			if ok {
				//fmt.Println("Rebalancing")
				err := v.Rebalance(g.ctx.K, g)
				if err != nil {
					log.Fatal(0)
				}

			} else {
				fmt.Println("Interface was supposedly NOT ok")
			}
		}
	}

	//fmt.Println("Running per bar")
	return g.ctx.strategy[0].PerBar(g.ctx.K, g)
}

func (k *GenerationK) Run() error {

	for k.ctx.K < k.ctx.length-1 {
		/*err := k.nextDate()
		if err != nil {
			log.Fatal("Could not increase the time")
		}*/

		err := k.nextGen()

		if err != nil {

			switch err {

			case EndOfBacktest:
				return err

			case FFToStartDate:
				continue

			case UnstablePeriod:
				continue

			case EndOfAsset:
				continue

			default:
				log.Print(err)

				return err

			}
		}

	}

	k.end()

	return nil
}

func (k *GenerationK) end() {
	fmt.Println("Running end")
	err := k.ctx.strategy[0].End(k.ctx.K, k)
	if err != nil {
		log.Fatal("Fatal end")
	}
	k.PrintRecords()
}

//AddDataManager is currently not used
func (k *GenerationK) SetDataManager() {}

//Returns an assets
func (k *GenerationK) GetAsset() *D.Asset {
	return k.ctx.GetAssets()[0]
}

//Returns an array of all assets
func (k *GenerationK) GetAssets() []*D.Asset {
	return k.ctx.GetAssets()
}

//GetAssetByName returns a pointer to the asset by that name
func (k *GenerationK) GetAssetByName(name string) *D.Asset {
	return k.ctx.GetAssetByName(name)
}

//SetComission is used to set the comission scheme is there is one
func (k *GenerationK) SetComission(comission Comission) {
	k.ctx.broker.SetComission(comission)
}

//AddAsset is used to add a pointer to an asset
func (k *GenerationK) AddAsset(asset *D.Asset) {
	k.ctx.AddAsset(asset)
}

//AddPortfolio is used to add a pointer to a portfolio to the backtest
func (k *GenerationK) SetPortfolio(portfolio *Portfolio) {
	k.ctx.portfolio = portfolio
	k.ctx.broker.portfolio = portfolio
}

//AddStrategy is used to add a strategy to the backtest
func (k *GenerationK) AddStrategy(strat Strategy) {
	k.ctx.SetStrategy(strat)
}

//SetBalance is used to set the balance when the backtest is started
func (k *GenerationK) SetBalance(balance float64) {
	k.ctx.portfolio.SetBalance(balance)
}

//AddStartDate is used to set the end date for the backtest
func (k *GenerationK) SetStartDate(startDate time.Time) {
	k.ctx.SetStartDate(startDate)
}

//AddEndDate is used to set the end date for the backtest
func (k *GenerationK) SetEndDate(endDate time.Time) {
	k.ctx.SetEndDate(endDate)
}

func (k *GenerationK) Sell(assetName string) (float64, error) {
	if asset, ok := k.ctx.assetMap[assetName]; ok {
		return k.sendOrder(k.ctx, SellOrder, MarketOrder, asset, k.ctx.datePointer, -1)
	}
	return 0, AssetDoesNotExist
}

//OrderSend is used to send an order to the broker, return an error if the asset does not exist
func (k *GenerationK) SendOrderFor(assetName string, direction Direction, orderType OrderType, qty int) (float64, error) {
	if asset, ok := k.ctx.assetMap[assetName]; ok {
		return k.sendOrder(k.ctx, direction, orderType, asset, k.ctx.datePointer, qty)
	}

	return 0, AssetDoesNotExist
}

//OrderSend is used to send an order to the broker, return an error if the asset does not exist
func (k *GenerationK) SendOrder(direction Direction, orderType OrderType, qty int) (float64, error) {
	if asset, ok := k.ctx.assetMap[k.ctx.asset.Name]; ok {
		return k.sendOrder(k.ctx, direction, orderType, asset, k.ctx.datePointer, qty)
	}

	return 0, AssetDoesNotExist
}

//orderSend is used to send an order to the broker
func (k *GenerationK) sendOrder(ctx *Context, direction Direction, orderType OrderType, asset *D.Asset, time time.Time, qty int) (float64, error) {
	orderStatus, _ := interface{}(ctx.strategy[0]).(OrderStatus)

	return ctx.broker.SendOrder(
		Order{
			direction: direction,
			orderType: orderType,
			Asset:     asset.Name,
			Price:     asset.Ohlc.Close[ctx.K-asset.AdjK],
			Time:      time,
			Qty:       qty,
		},
		orderStatus,
	)
}

//Assets returns an array of assets
func (k *GenerationK) Assets() []string {
	assets := make([]string, len(k.ctx.assets))
	for i, asset := range k.ctx.assets {
		assets[i] = asset.Name
	}

	return assets
}

func (k *GenerationK) Date() time.Time {
	return k.ctx.datePointer
}

/*func (k *GenerationK) SetUniverse(assets []string) {
	for i := range assets {
		k.AddAsset(NewAsset(assets[i], nil))
	}
}*/

//OwnPosition is used to find out if we have a holding in an asset
//and the assumption is that the strategy is using multiple assets
func (k *GenerationK) IsOwning(assetName string) (bool, error) {
	if _, ok := k.ctx.assetMap[assetName]; ok {
		return k.ctx.portfolio.IsOwning(assetName), nil
	}
	return false, AssetDoesNotExist
}

//Owning is used to find out if we have a holding and we are
//only processing 1 asset
func (k *GenerationK) Owning() bool {
	return k.ctx.portfolio.IsOwning(k.ctx.asset.Name)
}

type RecordStruct struct {
	Time     time.Time
	Variable map[string]string
}

func (k *GenerationK) Record(name string, record RecordStruct) {
	k.records[name] = append(k.records[name], record)
}

func (k *GenerationK) PositionSize(percentageOfCapital float64, price float64) int {
	return int(k.ctx.portfolio.cash * percentageOfCapital / price)
}

/*func (k *GenerationK) WriteRecords(fileName string) {
	file, err := os.Create(fileName + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for assetName, recordsArray := range k.records {

		for _, record := range recordsArray {

			err := writer.Write(value)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Fprintf(w, "\nAsset name %s\ttime %v\t\n", assetName, record.Time)

			for k, v := range record.Variable {
				fmt.Fprintf(w, "%s\t%s\taligned\t\n", k, v)
			}
		}

	}
}*/

func (k *GenerationK) PrintMatrix() {
	//records[time.Time]
}

func (k *GenerationK) PrintRecords() {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)

	for assetName, recordsArray := range k.records {

		for _, record := range recordsArray {
			fmt.Fprintf(w, "\nAsset name %s\ttime %v\t\n", assetName, record.Time)

			for k, v := range record.Variable {
				fmt.Fprintf(w, "%s\t%s\taligned\t\n", k, v)
			}
		}

		w.Flush()
	}
}
