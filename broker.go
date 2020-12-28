package generationk

import (
	log "github.com/sirupsen/logrus"
)

//Directon is used to describe an order
type Directon int

const (
	//BuyOrder order
	BuyOrder Directon = iota
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

type OrderType int

const (
	// A market order is an order to buy or sell a stock at the market’s
	// current best available price. A market order typically ensures
	// an execution but it does not guarantee a specified price.
	MarketOrder OrderType = iota

	//A limit order is an order to buy or sell a stock with a restriction on
	// the maximum price to be paid or the minimum price to be received
	// (the “limit price”). If the order is filled, it will only be at the
	// specified limit price or better. However, there is no assurance of
	// execution. A limit order may be appropriate when you think you can
	// buy at a price lower than—or sell at a price higher than—the
	// current quote.
	LimitOrder

	// A stop order is an order to buy or sell a stock at the market price once
	// the stock has traded at or through a specified price (the “stop price”).
	// If the stock reaches the stop price, the order becomes a market order and
	// is filled at the next available market price. If the stock fails to reach
	// the stop price, the order is not executed.
	StopOrder
)

//Broker is used to send orders
type Broker struct {
	portfolio *Portfolio
	//channel   chan Event
	callback  OrderStatus
	comission Comission
}

func (b *Broker) SetComission(comission Comission) {
	b.comission = comission
}

//SendOrder is used to place an order with the broker
func (b *Broker) SendOrder(order Order, orderstatus OrderStatus) {
	b.callback = orderstatus

	/*log.WithFields(log.Fields{
		"ordertype":                 order.orderType,
		"asset":                     order.Asset.name,
		"order.Asset.Ohlc[0].Close": order.Asset.ohlc[0].close,
		"time":                      order.Time,
		"amount":                    order.Amount,
		"qty":                       order.Qty,
	}).Info("BROKER>PLACE ORDER")*/

	switch order.direction {
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
	/*WithFields(log.Fields{
		"Order": order,
	}).Debug("BROKER> ACCEPTED")*/
	b.callback.OrderEvent(Accepted{})
}

func (b Broker) rejected(err error) {
	/*log.WithFields(log.Fields{
		"Error": err,
	}).Info("BROKER> REJECTED")*/
	b.callback.OrderEvent(Rejected{err: err})
}

func (b *Broker) buy(order Order) error {
	/*log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> BUY")*/
	amount := 0.0
	if order.Qty > ZERO {
		amount = getAmountForQty(order)
	}

	if order.Amount > EMPTY {
		amount = order.Amount
		qty := getQtyForAmount(order)
		order.Qty = qty
	}

	if b.comission != nil {
		amount += b.comission.GetComisson(order.Amount, order.Qty)
	} else {
		//log.Warning("Placing order without comission")
	}

	err := b.portfolio.subtractFromBalance(amount)
	if err != nil {
		b.rejected(err)

		return err
	}

	b.accepted(order)
	b.portfolio.AddHolding(Holding{
		qty:       order.Qty,
		assetName: order.Asset.name,
		price:     order.Asset.Close(),
		time:      order.Time,
	})
	//log.Debug("Calling the order event")

	b.callback.OrderEvent(Fill{
		Qty:       order.Qty,
		AssetName: order.Asset.name,
		Price:     order.Asset.Close(),
		Time:      order.Time,
	})

	/*log.Debug("Coming back after it")
	log.Debug("BROKER> Put FILL EVENT in queue")*/

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
