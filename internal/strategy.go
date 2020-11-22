package internal

type Strategy interface {
	Init(ctx *Context)
	Tick(ctx *Context)
}
