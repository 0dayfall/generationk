package internal

//Event type
type Event interface {
	Handle()
}

//Tick event type
type Tick struct{}

func (t Tick) Handle() {}

//Signal event type
type Signal struct{}

func (s Signal) Handle() {}

//Order event type
type Order struct{}

func (o Order) Handle() {}

//Fill event type
type Fill struct{}

func (f Fill) Handle() {}

//Data event type
type Data struct{}

func (d Data) Handle() {}

//Fill event type
type Quit struct{}

func (q Quit) Handle() {}
