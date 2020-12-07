package internal

//Updateable takes new data into account
type Updateable interface {
	Update(ctx *Context)
}

//Strategy needs to implement Orders in order to generate them
type Strategy interface {
	Setup(ctx *Context) error
	Tick(ctx *Context)
	OrderEvent(ctx *Context)
}
