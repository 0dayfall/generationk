package internal

//DataSource interface for getting data
type DataSource interface {
	GetData(period int)
	GetLatestData()
}
