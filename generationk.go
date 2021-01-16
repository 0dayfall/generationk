package generationk

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

//var o sync.Once

var AssetDoesNotExist = errors.New("Asset does not exist")
var FFToStartDate = errors.New("Fast forwarding to start date")
var EndOfBacktest = errors.New("End of backtest")
var EndOfData = errors.New("End of data")
var UnstablePeriod = errors.New("The stable period is not yet reached")
var Initialization = errors.New("Initialization in Once() failed")

//Callback is used in the strategy to give actions back to the backtest in progress
type Callback interface {
	Owning() bool
	IsOwning(assetName string) (bool, error)
	SendOrder(direction Direction, orderType OrderType, qty int) error
	SendOrderFor(assetName string, direction Direction, orderType OrderType, qty int) error
	Assets() []string
}

type GenerationK struct {
	ctx *Context
}

//NewGenerationK is used to create a new backtest
func NewGenerationK() *GenerationK {

	generationK := &GenerationK{
		ctx: NewContext(),
	}

	return generationK
}

func (g *GenerationK) inc() {
	g.ctx.K++
}

func (g *GenerationK) nextGen() error {
	defer g.inc()
	g.ctx.datePointer = g.ctx.asset.ohlc.Time[g.ctx.K]

	//Have to run this first so that we dont increase k by FF
	if g.ctx.K < 1 {
		//fmt.Printf("Once executed for %s\n\n", g.ctx.asset.name)
		err := g.ctx.strategy[0].Once(g.ctx, g.ctx.asset.ohlc)
		if err != nil {
			return err
		}
	}

	if g.ctx.asset.ohlc.Time[g.ctx.K].Before(g.ctx.startDate) {
		return FFToStartDate
	}

	if g.ctx.asset.ohlc.Time[g.ctx.K].After(g.ctx.endDate) {
		return EndOfBacktest
	}

	//Run setup after initperiod is finished
	if g.ctx.K < g.ctx.GetInitPeriod() {
		return UnstablePeriod
	}

	return g.ctx.strategy[0].PerBar(g.ctx.K, g)
}

func (k *GenerationK) Run() error {

	for k.ctx.K < k.ctx.length-1 {

		err := k.nextGen()
		if err != nil {
			switch err {
			case EndOfBacktest:
				return err
			case FFToStartDate:
				continue
			case UnstablePeriod:
				continue
			default:
				fmt.Println(err)
				return err
			}
		}

	}

	return nil
}

//AddDataManager is currently not used
func (k *GenerationK) SetDataManager() {}

//Returns an array of all assets
func (k *GenerationK) GetAsset() Asset {
	return k.ctx.GetAssets()[0]
}

//Returns an array of all assets
func (k *GenerationK) GetAssets() []Asset {
	return k.ctx.GetAssets()
}

//GetAssetByName returns a pointer to the asset by that name
func (k *GenerationK) GetAssetByName(name string) *Asset {
	return k.ctx.GetAssetByName(name)
}

//SetComission is used to set the comission scheme is there is one
func (k *GenerationK) SetComission(comission Comission) {
	k.ctx.broker.SetComission(comission)
}

//AddAsset is used to add a pointer to an asset
func (k *GenerationK) AddAsset(asset *Asset) {
	k.ctx.AddAsset(asset)
}

//AddPortfolio is used to add a pointer to a portfolio to the backtest
func (k *GenerationK) SetPortfolio(portfolio *Portfolio) {
	k.ctx.portfolio = portfolio
	k.ctx.broker.portfolio = portfolio
}

//AddStrategy is used to add a strategy to the backtest
func (k *GenerationK) AddStrategy(strat Strategy) {
	k.ctx.AddStrategy(strat)
}

//SetBalance is used to set the balance when the backtest is started
func (k *GenerationK) SetBalance(balance float64) {
	k.ctx.portfolio.SetBalance(balance)
}

//AddStartDate is used to set the end date for the backtest
func (k *GenerationK) SetStartDate(startDate time.Time) {
	k.ctx.AddStartDate(startDate)
}

//AddEndDate is used to set the end date for the backtest
func (k *GenerationK) SetEndDate(endDate time.Time) {
	k.ctx.AddEndDate(endDate)
}

//OrderSend is used to send an order to the broker, return an error if the asset does not exist
func (k *GenerationK) SendOrderFor(assetName string, direction Direction, orderType OrderType, qty int) error {
	if asset, ok := k.ctx.assetMap[k.ctx.asset.name]; ok {
		return k.sendOrder(k.ctx, direction, orderType, asset, k.ctx.datePointer, qty)
	}

	return AssetDoesNotExist
}

//OrderSend is used to send an order to the broker, return an error if the asset does not exist
func (k *GenerationK) SendOrder(direction Direction, orderType OrderType, qty int) error {
	if asset, ok := k.ctx.assetMap[k.ctx.asset.name]; ok {
		return k.sendOrder(k.ctx, direction, orderType, asset, k.ctx.datePointer, qty)
	}

	return AssetDoesNotExist
}

//orderSend is used to send an order to the broker
func (k *GenerationK) sendOrder(ctx *Context, direction Direction, orderType OrderType, asset *Asset, time time.Time, qty int) error {
	orderStatus, _ := interface{}(ctx.strategy[0]).(OrderStatus)

	err := ctx.broker.SendOrder(
		Order{
			direction: direction,
			orderType: orderType,
			Asset:     asset.name,
			Price:     asset.ohlc.Close[ctx.K],
			Time:      time,
			Qty:       qty,
		},
		orderStatus,
	)
	return err
}

//Assets returns an array of assets
func (k *GenerationK) Assets() []string {
	assets := make([]string, len(k.ctx.assets))
	for i, asset := range k.ctx.assets {
		assets[i] = asset.name
	}

	return assets
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
	return k.ctx.portfolio.IsOwning(k.ctx.asset.name)
}

//min returns the smaller of x or y.
func min(x, y int) int {
	if x > y {
		return y
	}

	return x
}

type Job interface {
	GetId() int
	GetFileName() string
	GetResult() float64
	GetParams() []float64
}

type JobStruct struct {
	Id       int
	FileName string
	Result   float64
	Params   []Params
}

func (j *JobStruct) GetId() int {
	return j.Id
}

func (j *JobStruct) GetFileName() string {
	return j.FileName
}

func (j *JobStruct) GetResult() float64 {
	return j.Result
}

func (j *JobStruct) SetParams(params ...Params) {
	j.Params = params
}

func produce(folderPath string, jobs chan<- *JobStruct) {
	// Generate jobs:
	files, err := filepath.Glob(folderPath + "*.csv")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("files %s", files)
	id := 0

	for _, fileName := range files {
		id++
		//fmt.Printf("ADDING JOB FOR %s", fileName)
		jobs <- &JobStruct{
			Id:       id,
			FileName: fileName,
			Result:   0,
			Params:   []Params{{Low: 10.0, High: 30.0}, {Low: -30.0, High: -20.0}},
		}
	}
	//fmt.Println("CLOSING JOBS")
	close(jobs)
}

func consume(id int, strategy Strategy, jobs <-chan *JobStruct, results chan<- *JobStruct, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {

		//fmt.Printf("READING JOBS %v", job)

		//Perform work
		dataManager := NewCSVDataManager()
		asset, err := dataManager.ReadCSVFile(job.FileName)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(0)
		}

		//Clunky way to check parameters
		//cross := new(RMICrossStrategy)
		val := reflect.ValueOf(strategy)
		if val.Kind() == reflect.Ptr {
			val = reflect.Indirect(val)
		}
		newThing := reflect.New(val.Type()).Interface().(Strategy)

		params := newThing.GetParams()

		fmt.Printf("%f %f %f\n", params[0].High, params[0].Low, params[0].Value)
		fmt.Printf("%f %f %f\n\n", params[1].High, params[1].Low, params[1].Value)
		fmt.Printf("SIZE: %d\n\n", len(params))

		for i := 0; i < len(params); i++ {
			for s := params[0].Low; s < params[0].High; s++ {
				for k := params[1].Low; k < params[1].High; k++ {
					//fmt.Printf("PARAMS VÄRDE I LOOPEN: %v", params[0].GetValue())
					params[0].Value = s
					params[1].Value = k
					portfolio := NewPortfolio()
					portfolio.SetBalance(100000)
					genk := NewGenerationK()
					genk.SetPortfolio(portfolio)
					/*				cross.SetParams(K.Params{Value: float64(i)})
									cross.SetSellParams(K.Params{Value: float64(s)})*/

					//fmt.Printf("PARAMS VÄRDE I LOOPEN: %v", strategy)
					genk.AddStrategy(newThing)

					now := time.Now()
					start := now.AddDate(-15, -9, -2)
					genk.SetStartDate(start)
					now = time.Now()
					end := now.AddDate(0, -3, -2)
					genk.SetEndDate(end)
					genk.AddAsset(asset)

					runErr := genk.Run()
					if runErr != nil {
						fmt.Print(runErr.Error())
						os.Exit(0)
					}
					/*params[0].Value = s
					params[0].Value = k*/

					//Fill out the result
					job.Result = portfolio.GetBalance()
					job.SetParams(*params[0])
					//Send the results back
					results <- job
				}
			}
		}
	}
}

/*func run(localFilename string, wg *sync.WaitGroup, c chan float64) {
	defer wg.Done()

	}

}*/

func analyze(results <-chan *JobStruct, wg2 *sync.WaitGroup) {
	defer wg2.Done()

	f, err := os.Create("data.txt")
	w := bufio.NewWriter(f)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer f.Close()
	var m = make(map[string]JobStruct)

	for job := range results {
		_, err2 := fmt.Fprintf(w, "id: %d, filename: %s, result: %f\n", job.Id, job.FileName, job.Result)
		if m[job.FileName].Result < job.Result {
			m[job.FileName] = *job
		}

		if err2 != nil {
			fmt.Print(err2.Error())
		}
	}

	for _, v := range m {
		fmt.Printf("Param: %s\n", v.FileName)
		fmt.Printf("Param: %f\n", v.Params[0].Value)
		//fmt.Printf("Param: %f\n", v.Params[1].Value)
		fmt.Printf("Result %f\n\n", v.Result)
	}
	//fmt.Printf("M: %v", m)
}

func Run(folderPath string, strategy Strategy) {
	//defer profile.Start().Stop()
	//t.Parallel()
	jobs := make(chan *JobStruct, 100)    // Buffered channel
	results := make(chan *JobStruct, 100) // Buffered channel
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	// Start consumers:
	for i := 0; i < 5; i++ { // 5 consumers
		wg.Add(1)
		go consume(i, strategy, jobs, results, &wg)
	}

	go produce(folderPath, jobs)

	wg2.Add(1)
	go analyze(results, &wg2)

	wg.Wait()
	close(results)

	wg2.Wait()
}

func RunStrategyOnAssets(strategy Strategy, folderPath string) {
	files, err := filepath.Glob(folderPath + "*.csv")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("files %s\n\n\n", files)
	//d.ReadCSVFilesAsync(files)
	portfolio := NewPortfolio()
	portfolio.SetBalance(100000)

	var wg sync.WaitGroup

	y := 0

	for _, fileName := range files {
		wg.Add(1)
		go func(localFilename string, strategy Strategy) {
			genk := NewGenerationK()
			genk.SetPortfolio(portfolio)
			genk.AddStrategy(strategy)

			now := time.Now()
			start := now.AddDate(-15, -9, -2)
			genk.SetStartDate(start)
			now = time.Now()
			end := now.AddDate(0, -3, -2)
			genk.SetEndDate(end)

			//genk.RunEventBased()
			dataManager := NewCSVDataManager()
			//dataManager.SetHandler(genk)
			//genk.AddDataManager(dataManager)

			//dataManager.ReadCSVFilesAsync([]string{"test/data/ABB.csv", "test/data/ASSAb.csv"})
			asset, err := dataManager.ReadCSVFile(localFilename)
			if err != nil {
				fmt.Print(err.Error())
				os.Exit(0)
			}
			genk.AddAsset(asset)

			runErr := genk.Run()
			fmt.Print(runErr.Error())

			wg.Done()
		}(fileName, strategy)
		y++
	}
	wg.Wait()
}
