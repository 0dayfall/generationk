package generationk

import (
	"time"
)

type Statistics struct {
	records map[string][]Statistic
	cash    float64
}

func (stat *Statistic) New(cash float64) *Statistics {
	return &Statistics{
		records: map[string][]Statistic{},
		cash:    cash, 
	}
}

func (stat *Statistics) record(stockName string, statistic Statistic) {
	list := stat.records[stockName]
	list = append(list, statistic)
	stat.records[stockName] = list
}

func (stat *Statistics) profits(stockName string) float64 {
	//Assumed list is sorted

}

//This is what we are owning, a holding
type Statistic struct {
	qty       int
	buy       bool
	assetName string
	price     float64
	time      time.Time
}
