package generationk

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	D "github.com/0dayfall/generationk/data"
)

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

func produce(ctx *Context, dm *D.DataManager, jobs chan<- *JobStruct) {
	// Generate jobs:
	files, err := filepath.Glob(ctx.dataPath + "*.csv")
	if err != nil {
		log.Fatal(err)
	}

	id := 0
	//fmt.Println(files)
	for _, fileName := range files {
		id++
		//fmt.Println(fileName)
		jobs <- &JobStruct{
			Id:       id,
			FileName: fileName,
			Result:   0,
			Params:   []Params{{Low: 10.0, High: 30.0}, {Low: -30.0, High: -20.0}},
		}
	}

	close(jobs)
}

func consume(id int, ctx *Context, dm *D.DataManager, jobs <-chan *JobStruct, results chan<- *JobStruct, wg *sync.WaitGroup) {
	defer wg.Done()

	//fmt.Printf("Ranging jobs")

	for job := range jobs {

		//fmt.Printf("READING JOBS %v", job)

		//Perform work
		asset, err := dm.ReadCSVFile(job.FileName)
		if err != nil {
			log.Fatal(err)
		}

		//Clunky way to check parameters
		//cross := new(RebalanceStrategy)
		val := reflect.ValueOf(ctx.GetStrategy())
		if val.Kind() == reflect.Ptr {
			val = reflect.Indirect(val)
		}
		newThing := reflect.New(val.Type()).Interface().(Strategy)

		params := newThing.GetParams()

		//fmt.Printf("%f %f %f\n", params[0].High, params[0].Low, params[0].Value)
		//fmt.Printf("%f %f %f\n\n", params[1].High, params[1].Low, params[1].Value)
		//		fmt.Printf("SIZE: %d\n\n", len(params))

		for i := 0; i < len(params); i++ {
			for s := params[0].Low; s < params[0].High; s += 1 {
				for k := params[1].Low; k < params[1].High; k += 1 {
					//fmt.Printf("PARAMS VÄRDE I LOOPEN: %v", params[0].GetValue())
					params[0].Value = s
					params[1].Value = k
					//fmt.Printf("PARAMS VÄRDE I LOOPEN: %v", params[0].GetValue())
					//fmt.Printf("PARAMS VÄRDE I LOOPEN: %v", params[1].GetValue())
					portfolio := NewPortfolio(1000000)
					//portfolio.SetBalance(100000)
					genk := NewGenerationK()
					genk.SetPortfolio(portfolio)
					//fmt.Print(portfolio)

					//fmt.Printf("PARAMS VÄRDE I LOOPEN:")
					genk.AddStrategy(newThing)
					genk.SetStartDate(ctx.GetStartDate())
					genk.SetEndDate(ctx.GetEndDate())
					genk.AddAsset(asset)

					err := genk.Run()
					if err != nil {
						log.Fatal(err)
					}

					//Fill out the result
					//fmt.Print(portfolio)
					job.Result = portfolio.GetBalance()
					//fmt.Printf("====>%v", job.Result)
					//fmt.Println(portfolio.GetBalance())
					job.SetParams(*params[0], *params[1])
					//Send the results back
					results <- job
				}
			}
		}
	}
}

func analyze(results <-chan *JobStruct, wg2 *sync.WaitGroup) {
	fmt.Println("Analyzing results")
	defer wg2.Done()

	f, err := os.Create("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	var m = make(map[string]JobStruct)
	fmt.Print("Waiting for results!!!\n")
	for job := range results {
		fmt.Printf("id: %d, filename: %s, result: %f\n", job.Id, job.FileName, job.Result)
		_, err := fmt.Fprintf(w, "id: %d, filename: %s, result: %f\n", job.Id, job.FileName, job.Result)
		if err != nil {
			log.Fatal(err)
		}

		if m[job.FileName].Result < job.Result {
			m[job.FileName] = *job
		}
	}

	for _, v := range m {
		fmt.Printf("Param: %s\n", v.FileName)
		for _, param := range v.Params {
			fmt.Printf("Param: %f\n", param.Value)
		}
		fmt.Printf("Result %f\n\n", v.Result)
	}

}

func RunAsJobs(ctx *Context, dm *D.DataManager) {
	//defer profile.Start().Stop()
	//t.Parallel()
	jobs := make(chan *JobStruct, 100)    // Buffered channel
	results := make(chan *JobStruct, 100) // Buffered channel
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	// Start consumers:
	for i := 0; i < 5; i++ { // 5 consumers
		wg.Add(1)
		go consume(i, ctx, dm, jobs, results, &wg)
	}

	go produce(ctx, dm, jobs)

	wg2.Add(1)
	go analyze(results, &wg2)

	wg.Wait()
	close(results)

	wg2.Wait()
}

func RunPlain(ctx *Context, dm *D.DataManager) {
	files, err := filepath.Glob(filepath.Clean(ctx.dataPath) + string(os.PathSeparator) + "*.csv")

	if err == filepath.ErrBadPattern {
		log.Fatal(err)
	}

	portfolio := NewPortfolio(100000)
	//portfolio.SetBalance(100000)

	for _, fileName := range files {
		genk := NewGenerationK()
		genk.SetPortfolio(portfolio)
		genk.AddStrategy(ctx.GetStrategy())

		genk.SetStartDate(ctx.GetStartDate())
		genk.SetEndDate(ctx.GetEndDate())

		asset, err := dm.ReadCSVFile(fileName)
		if err != nil {
			log.Fatal(err)
		}
		genk.AddAsset(asset)

		runErr := genk.Run()
		if runErr != nil {
			log.Fatal(runErr)
		}
	}
}

func RunParallell(ctx *Context, dm *D.DataManager, cash float64) {
	files, err := filepath.Glob(filepath.Clean(ctx.dataPath) + string(os.PathSeparator) + "*.csv")

	if err == filepath.ErrBadPattern {
		log.Fatal(err)
	}

	portfolio := NewPortfolio(cash)
	//portfolio.SetBalance(500000)

	genk := NewGenerationK()
	genk.SetPortfolio(portfolio)
	genk.AddStrategy(ctx.GetStrategy())

	genk.SetStartDate(ctx.GetStartDate())
	genk.SetEndDate(ctx.GetEndDate())

	for _, fileName := range files {

		asset, err := dm.ReadCSVFile(fileName)
		if err != nil {
			fmt.Printf("Filename> %s %s", fileName, err.Error())
			os.Exit(1)
		}

		genk.AddAsset(asset)

	}

	dm.CreatePadding(genk.GetAssets())

	runErr := genk.Run()
	if runErr != nil {
		fmt.Printf("Error, portfolio value: %f", portfolio.GetBalance())
		log.Fatal(runErr)
	}

	fmt.Printf("Portfolio value: %f", portfolio.GetBalance())
}

func RunStrategyOnAssets(ctx *Context, dm *D.DataManager) {
	files, err := filepath.Glob(ctx.dataPath + "*.csv")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("files %s\n\n\n", files)
	//d.ReadCSVFilesAsync(files)
	portfolio := NewPortfolio(100000)
	//portfolio.SetBalance(100000)

	var wg sync.WaitGroup

	y := 0

	for _, fileName := range files {
		wg.Add(1)
		go func(localFilename string, ctx *Context, dm *D.DataManager) {
			genk := NewGenerationK()
			genk.SetPortfolio(portfolio)
			genk.AddStrategy(ctx.GetStrategy())

			genk.SetStartDate(ctx.GetStartDate())
			genk.SetEndDate(ctx.GetEndDate())

			asset, err := dm.ReadCSVFile(localFilename)
			if err != nil {
				log.Fatal(err)
			}
			genk.AddAsset(asset)

			runErr := genk.Run()
			if runErr != nil {
				log.Fatal(runErr)
			}

			wg.Done()
		}(fileName, ctx, dm)
		y++
	}
	wg.Wait()
}
