package internal

import (
	"time"

	log "github.com/sirupsen/logrus"
)

//Event type
type Event interface {
	Handle()
}

type DataEvent struct{}

func (d DataEvent) Handle() {}

type Order struct {
	Ordertype OrderType
	Asset     *Asset
	Time      time.Time
	Amount    float64
}

func (o Order) Handle() {}

type Fill struct {
	Qty       int
	AssetName string
	Time      time.Time
}

func (f Fill) Handle() {}

func (f Fill) String() {
	log.WithFields(log.Fields{
		"Qty":       f.Qty,
		"AssetName": f.AssetName,
		"Time":      f.Time,
	}).Debug("Fill> ")
}

//Tick event type
type Tick struct{}

func (t Tick) Handle() {}

//Signal event type
type Signal struct{}

func (s Signal) Handle() {}

//Data event type
type Data struct{}

func (d Data) Handle() {}

//Fill event type
type Quit struct{}

func (q Quit) Handle() {}
