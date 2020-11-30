package internal

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Direction int

const (
	Long Direction = iota
	Short
)

type Position struct {
	amount    float64
	assetName string
	time      time.Time
	price     float64
	comission float64
}

type Portfolio struct {
	Positions []Position
	cash      float64
	channel   chan Event
}

func (p *Portfolio) Add(position Position) {
	log.WithFields(log.Fields{
		"asset": position.assetName,
		"time":  position.time,
		"price": position.price,
	}).Debug("Adding position to portfolio")
	if p.Positions != nil {
		p.Positions = append(p.Positions, position)
	}
}

func (p *Portfolio) SetCash(amount float64) {
	p.cash = amount
}

func (p *Portfolio) GetCash() float64 {
	return p.cash
}
