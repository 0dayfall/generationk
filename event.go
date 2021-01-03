package generationk

import (
	"fmt"
	"time"
)

//DataHandler is the interface used to recieve data
//from any data producing function. It can be used to feed data to generationK
type DataHandler interface {
	DataEvent(dataEvent Event)
}

//OrderStatus is a callback interface used to recieve
//information about orders, it is used by the broker
type OrderStatus interface {
	OrderEvent(orderEvent Event)
}

//Event type
type Event interface {
	String() string
}

//DataEvent is a data structure used to carry OHLC data
type DataEvent struct {
	Name string
	Ohlc OHLC
}

func (d DataEvent) String() string {
	return fmt.Sprintf("$DATAEVENT %s", d.Name)
}

//Order describes an order that is used to buy / sell an asset
type Order struct {
	direction Direction
	orderType OrderType
	Asset     *Asset
	Time      time.Time
	Amount    float64
	Qty       int
}

func (o Order) String() string {
	return fmt.Sprintf("$ORDER %v %v %v %f %d", o.orderType, o.Asset, o.Time, o.Amount, o.Qty)
}

//Accepted is a status of the order to indicate that an order has been accepted by the broker.
type Accepted struct {
}

func (a Accepted) String() string {
	return "$ACCEPTED"
}

//Submitted is a status used after an order is to be processed by the broker
type Submitted struct {
}

func (s Submitted) String() string {
	return "$SUBMITTED"
}

//PartialFill is used to giv enotice to the strategy that a partial fill took place
type PartialFill struct {
}

func (pf PartialFill) String() string {
	return "$PARTIALFILL"
}

//Fill is used to indicate to the implementer of OrderStatus that an
//order has been filled
type Fill struct {
	Qty       int
	Price     float64
	AssetName string
	Time      time.Time
}

func (f Fill) String() string {
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
