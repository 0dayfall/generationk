package generationk

import (
	"time"

	log "github.com/sirupsen/logrus"
)

//Event type
type Event interface {
	Handle()
}

//DataEvent is for sending data
type DataEvent struct {
	Name string
	Ohlc OHLC
}

//Handle iM not sure what it si used for
func (d DataEvent) Handle() {
	log.WithFields(log.Fields{
		"Name": d.Name,
		"Ohlc": d.Ohlc,
	}).Debug("DataEvent$ ")
}

//Order describes an order
type Order struct {
	Ordertype OrderType
	Asset     *Asset
	Time      time.Time
	Amount    float64
	Qty       int
}

func (o Order) Handle() {}

type Fill struct {
	Qty       int
	Price     float64
	AssetName string
	Time      time.Time
}

func (f Fill) Handle() {}

func (f Fill) String() {
	log.WithFields(log.Fields{
		"Qty":       f.Qty,
		"Price":     f.Price,
		"AssetName": f.AssetName,
		"Time":      f.Time,
	}).Debug("Fill$")
}

type Rejected struct {
	message string
}

func (r Rejected) Handle() {}

func (r Rejected) String() {
	log.WithFields(log.Fields{
		"Message": r.message,
	}).Debug("Fill$")
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

//Quit event type
type Quit struct{}

func (q Quit) Handle() {}
