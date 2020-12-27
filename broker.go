package generationk

import (
	log "github.com/sirupsen/logrus"
)

//OrderType is used to describe an order
type OrderType int

const (
	//BuyOrder order
	BuyOrder OrderType = iota
	//SellOrder order
	SellOrder
	//ShortOrder order
	ShortOrder
	//CoverOrder short order
	CoverOrder
	//Zero = 0
	ZERO = 0
	//Empty account
	EMPTY = 0.0
)

//Broker is used to send orders
type Broker struct {
	portfolio *Portfolio
	//channel   chan Event
	callback OrderStatus
}

//SendOrder is used to place an order with the broker
func (b *Broker) SendOrder(order Order, orderstatus OrderStatus) {
	b.callback = orderstatus

	log.WithFields(log.Fields{
		"ordertype":                 order.Ordertype,
		"asset":                     order.Asset.name,
		"order.Asset.Ohlc[0].Close": order.Asset.ohlc[0].close,
		"time":                      order.Time,
		"amount":                    order.Amount,
		"qty":                       order.Qty,
	}).Info("BROKER>PLACE ORDER")

	switch order.Ordertype {
	case BuyOrder:
		err := b.buy(order)
		if err != nil {
			b.rejected(err)
		}
	case SellOrder:
		b.sell(order)
	case ShortOrder:
		b.sellshort(order)
	case CoverOrder:
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
	}).Debug("BROKER> ACCEPTED")
	b.callback.OrderEvent(Accepted{})
}

func (b Broker) rejected(err error) {
	log.WithFields(log.Fields{
		"Error": err,
	}).Info("BROKER> REJECTED")
	b.callback.OrderEvent(Rejected{err: err})
}

func (b *Broker) buy(order Order) error {
	/*log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> BUY")*/
	if order.Qty > ZERO {
		err := b.portfolio.subtractFromBalance(getAmountForQty(order))
		if err != nil {
			b.rejected(err)

			return err
		}
		qty := getQtyForAmount(order)
		order.Qty = qty
	}

	if order.Amount > EMPTY {
		err := b.portfolio.subtractFromBalance(order.Amount)
		if err != nil {
			b.rejected(err)

			return err
		}

		qty := getQtyForAmount(order)
		order.Qty = qty
	}

	b.accepted(order)
	b.portfolio.AddHolding(Holding{
		qty:       order.Qty,
		assetName: order.Asset.name,
		price:     order.Asset.Close(),
		time:      order.Time,
	})
	log.Debug("Calling the order event")

	b.callback.OrderEvent(Fill{
		Qty:       order.Qty,
		AssetName: order.Asset.name,
		Price:     order.Asset.Close(),
		Time:      order.Time,
	})

	log.Debug("Coming back after it")
	log.Debug("BROKER> Put FILL EVENT in queue")

	return nil
}

func (b *Broker) sell(order Order) {
	log.WithFields(log.Fields{
		"Order": order.Asset.name,
	}).Info("BROKER> SELL")

	if order.Qty > ZERO {
		b.portfolio.addToBalance(getAmountForQty(order))
	}

	if order.Amount > EMPTY {
		b.portfolio.addToBalance(order.Amount)
	}

	b.accepted(order)

	b.portfolio.RemoveHolding(Holding{
		qty:       order.Qty,
		assetName: order.Asset.name,
		price:     order.Asset.Close(),
		time:      order.Time,
	})

	b.callback.OrderEvent(Fill{
		Qty: -order.Qty, AssetName: order.Asset.name,
		Price: order.Asset.Close(),
		Time:  order.Time,
	})

	log.Debug("BROKER> Put FILL EVENT in queue")
}

func (b *Broker) sellshort(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Debug("BROKER> SELLSHORT")
}

func (b *Broker) cover(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Debug("BROKER> COVER")
}
