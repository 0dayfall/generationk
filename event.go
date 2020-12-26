package generationk

import (
	"fmt"
	"time"
)

type DataHandler interface {
	DataEvent(dataEvent Event)
}

//OrderStatus is a callback interface used to recieve information about orders, it is used by the broker
type OrderStatus interface {
	OrderEvent(orderEvent Event)
}

//Event type
type Event interface {
	String() string
}

//DataEvent is for sending data
type DataEvent struct {
	Name string
	Ohlc OHLC
}

//Handle iM not sure what it si used for
func (d DataEvent) String() string {
	/*log.WithFields(log.Fields{
		"Name": d.Name,
		"Ohlc": d.Ohlc,
	}).Debug("DataEvent$ ")*/
	return fmt.Sprintf("$DATAEVENT %s", d.Name)
}

//Order describes an order
type Order struct {
	Ordertype OrderType
	Asset     *Asset
	Time      time.Time
	Amount    float64
	Qty       int
}

func (o Order) String() string {
	return fmt.Sprintf("$ORDER %v %v %v %f %d", o.Ordertype, o.Asset, o.Time, o.Amount, o.Qty)
}

//Accepted is a status of the order
type Accepted struct {
}

func (a Accepted) String() string {
	return "$ACCEPTED"
}

type Submitted struct {
}

func (s Submitted) String() string {
	return "$SUBMITTED"
}

type PartialFill struct {
}

func (pf PartialFill) String() string {
	return "$PARTIALFILL"
}

type Fill struct {
	Qty       int
	Price     float64
	AssetName string
	Time      time.Time
}

func (f Fill) String() string {
	/*log.WithFields(log.Fields{
		"Qty":       f.Qty,
		"Price":     f.Price,
		"AssetName": f.AssetName,
		"Time":      f.Time,
	}).Debug("Fill$")*/
	return fmt.Sprintf("%d %f %s %v", f.Qty, f.Price, f.AssetName, f.Time)
}

//Rejected type is for order that can not be executed
type Rejected struct {
	err error
}

func (r Rejected) String() string {
	/*log.WithFields(log.Fields{
		"Message": r.err.Error(),
	}).Debug("REJECTED$")*/
	return r.err.Error()
}

//Tick event type
type Tick struct{}

func (t Tick) String() string {
	return "$TICK"
}

//Quit event type
type Quit struct{}

func (q Quit) String() string {
	return "$QUIT"
}
