package generationk

//Updateable takes new data into account
/*type Updateable interface {
	Update([]float64)
}*/

/*type SetupStrategy interface {
	AddIndicator(indicator indicators.Indicator)
	SetInitPeriod(period int)
}*/

//Strategy needs to implement Orders in order to generate them
type Strategy interface {
	Setup(ctx *Context) error
	Tick(genk GenkCallback)
}
