package internal

import "time"

type OrderNotfication interface {
	OrderConfirmation() Position
}

type OrderType int

const (
	Buy OrderType = iota
	Sell
	SellShort
	Cover
)

type Broker struct {
	notifiers []OrderNotfication
	portfolio Portfolio
}

func (b *Broker) Buy(asset *Asset, time time.Time, amount float64) {
	pos := &Position{
		amount:    amount,
		assetName: asset.Name,
		time:      time,
		price:     asset.Close(),
		comission: 0,
	}
	b.portfolio.Add(*pos)

	/*for notify := range b.notifiers {
		notify.OrderConfirmaion()
	}*/
}

func (b *Broker) Sell(asset *Asset, amount int) {

}
