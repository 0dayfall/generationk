package generationk

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
	callback  OrderStatus
	comission Comission
}

//SetComission is used to set a comission scheme
func (b *Broker) SetComission(comission Comission) {
	b.comission = comission
}

//SendOrder is used to place an order with the broker
func (b *Broker) SendOrder(order Order, orderstatus OrderStatus) {
	b.callback = orderstatus

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

//getAmountForQty is used to get the amount needed to buy a certain quantity
func getAmountForQty(order Order) float64 {
	return order.Asset.Close() * float64(order.Qty)
}

//getQtyForAmount is used to get how many many stocks we get for a certain amount
func getQtyForAmount(order Order) int {
	return int(order.Amount / order.Asset.Close())
}

//accpeted is used to send an order event with an accepted event
func (b Broker) accepted(order Order) {
	b.callback.OrderEvent(Accepted{})
}

//rejected is used to send a rejected order event
func (b Broker) rejected(err error) {
	b.callback.OrderEvent(Rejected{err: err})
}

//buy is used to buy a qty or a possible qty from an amount, if a comission is
// set is will be used and deducted from the account
func (b *Broker) buy(order Order) error {
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

	b.callback.OrderEvent(Fill{
		Qty:       order.Qty,
		AssetName: order.Asset.name,
		Price:     order.Asset.Close(),
		Time:      order.Time,
	})

	return nil
}

//sell is used to sell a holding and book and the profits or losses
func (b *Broker) sell(order Order) {

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

}

//sellshort is not implemented
func (b *Broker) sellshort(order Order) {
	//TODO: Implement
}

//cover is used to cover a short, not implemented
func (b *Broker) cover(order Order) {
	//TODO: Implement
}
