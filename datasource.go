package generationk

//DataSource interface for getting data, not used in this implementation
type DataSource interface {
	GetData(period int)
	GetLatestData()
}
