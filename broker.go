package generationk

import (
	log "github.com/sirupsen/logrus"
)

//OrderType is used to describe an order
type OrderType int

const (
	//Buy order
	Buy OrderType = iota
	//Sell order
	Sell
	//SellShort order
	SellShort
	//Cover short order
	Cover
)

//OrderStatus is a callback interface used to recieve information about orders
type OrderStatus interface {
	OrderEvent(orderEvent Event)
}

//Broker is used to send orders
type Broker struct {
	portfolio Portfolio
	channel   chan Event
	callback  OrderStatus
}

//PlaceOrder is used to place an order with the broker
func (b *Broker) PlaceOrder(order Order, orderstatus OrderStatus) {
	b.callback = orderstatus

	log.WithFields(log.Fields{
		"ordertype": order.Ordertype,
		"asset":     (*order.Asset).Name,
		"time":      order.Time,
		"amount":    order.Amount,
	}).Debug("BROKER>PLACE BUY ORDER")

	switch order.Ordertype {
	case Buy:
		b.buy(order)
	case Sell:
		b.sell(order)
	case SellShort:
		b.sellshort(order)
	case Cover:
		b.cover(order)
	}
}

func getAmountForQty(order Order) float64 {
	return order.Asset.Close() * float64(order.Qty)
}

func getQtyForAmount(order Order) int {
	return int(order.Amount / order.Asset.Close())
}

func (b Broker) accepted(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> ACCEPTED")
	b.callback.OrderEvent(Accepted{})
}

func (b Broker) rejected(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> REJECTED")
	b.callback.OrderEvent(Rejected{message: "Insufficient funds"})
}

func (b *Broker) buy(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> BUY")
	if order.Qty > 0 {
		err := b.portfolio.updateCash(getAmountForQty(order))
		if err != nil {
			b.rejected(order)
			return
		}
	}
	if order.Amount > 0.0 {
		err := b.portfolio.updateCash(order.Amount)
		if err != nil {
			b.rejected(order)
			return
		}
		qty := getQtyForAmount(order)
		order.Qty = qty
	}

	b.accepted(order)
	b.portfolio.AddHolding(Holding{Qty: order.Qty, AssetName: order.Asset.Name, Price: order.Asset.Close(), Time: order.Time})
	log.Info("Calling the order event")
	b.callback.OrderEvent(Fill{Qty: order.Qty, AssetName: order.Asset.Name, Price: order.Asset.Close(), Time: order.Time})
	log.Info("Coming back after it")
	log.Info("BROKER> Put FILL EVENT in queue")
}

func (b *Broker) sell(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> SELL")
	b.callback.OrderEvent(Fill{Qty: order.Qty, AssetName: order.Asset.Name, Price: order.Asset.Close(), Time: order.Time})
	log.Info("BROKER> Put FILL EVENT in queue")
}

func (b *Broker) sellshort(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> SELLSHORT")
}

func (b *Broker) cover(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> COVER")
}
