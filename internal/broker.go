package internal

type OrderNotfication interface {
	OrderConfirmation() Position
}

type Broker struct {
	notifiers []OrderNotfication
	cash      float64
}

func (b *Broker) Buy(asset Asset, amount int) {
	pos := &Position{}
	for notify := range b.notifiers {
		notify.OrderConfirmaion()
	}
}

func (b *Broker) Sell(asset Asset, amount int) {

}
