package datasource

//DataInterface can be from CVS or from real time
type DataInterface interface {
	getData(period int) []float64
	getLatestData() float64
}

//RealTimeData is used as an interface for real time data
type RealTimeData interface {
}

//HistoricData is used as an interface for historic data
type HistoricData interface {
}
