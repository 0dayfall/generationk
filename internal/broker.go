package internal

import (
	"time"

	log "github.com/sirupsen/logrus"
)

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
	channel   chan Event
}

func (b *Broker) Order(ordertype OrderType, asset *Asset, time time.Time, amount float64) {
	log.WithFields(log.Fields{
		"ordertype": ordertype,
		"asset":     asset.Name,
		"time":      time,
		"amount":    amount,
	}).Debug("Creating order")
	if ordertype == Buy {
		b.buy(asset, time, amount)
	}
}

func (b *Broker) buy(asset *Asset, time time.Time, amount float64) {
	pos := &Position{
		amount:    amount,
		assetName: asset.Name,
		time:      time,
		price:     asset.Close(),
		comission: 0,
	}
	b.portfolio.Add(*pos)
	log.Debug("Sending fill event to channel")
	go func() {
		b.channel <- Fill{}
		for notify := range b.notifiers {
			notify.OrderConfirmaion()
		}
	}()
}

func (b *Broker) Sell(asset *Asset, amount int) {

}
