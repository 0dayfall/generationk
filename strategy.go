package generationk

//Updateable takes new data into account
type Updateable interface {
	Update([]float64)
}

//Strategy needs to implement Orders in order to generate them
type Strategy interface {
	Setup(ctx *Context) error
	Tick(genk GenkCallback)
}
