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

type Portfolio struct {
	Holdings []Holding
	cash     float64
}

type Holding struct {
	Qty       int
	AssetName string
	Time      time.Time
}

func (p *Portfolio) Fill(fill Fill) {
	p.AddHolding(Holding{Qty: fill.Qty, AssetName: fill.AssetName, Time: fill.Time})
}

func (p *Portfolio) AddHolding(position Holding) {
	log.WithFields(log.Fields{
		"asset": position.AssetName,
		"time":  position.Time,
		"Qty":   position.Qty,
	}).Debug("PORTFOLIO> Adding position to portfolio")
	if p.Holdings != nil {
		p.Holdings = append(p.Holdings, position)
	}
}

func (p *Portfolio) SetCash(amount float64) {
	p.cash = amount
}

func (p *Portfolio) GetCash() float64 {
	return p.cash
}
