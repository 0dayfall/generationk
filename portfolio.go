package generationk

import (
	"errors"
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
	Price     float64
	Time      time.Time
}

//IsOwning is used to find out if a position is already owned in this asset
func (p Portfolio) IsOwning(assetName string) bool {
	for k := range p.Holdings {
		if p.Holdings[k].AssetName == assetName {
			return true
		}
	}
	return false
}

func (p *Portfolio) Fill(fill Fill) {

	go p.AddHolding(Holding{Qty: fill.Qty, Price: fill.Price, AssetName: fill.AssetName, Time: fill.Time})
}

func (p *Portfolio) SellHolding(position Holding) {
	log.WithFields(log.Fields{
		"asset": position.AssetName,
		"time":  position.Time,
		"price": position.Price,
		"qty":   position.Qty,
	}).Debug("PORTFOLIO> Adding position to portfolio")
	if p.Holdings != nil {
		p.Holdings = append(p.Holdings, position)
	}
	p.updateCash(+float64(position.Qty) * position.Price)
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
	p.updateCash(-float64(position.Qty) * position.Price)
}

func (p *Portfolio) updateCash(cost float64) error {
	balance := p.cash + cost
	if balance < 0 {
		return errors.New("Balance < 0")
	}
	p.cash = balance
	log.WithFields(log.Fields{
		"Balance": p.cash,
	}).Debug("PORTFOLIO> Updating balance")
	return nil
}

func (p *Portfolio) SetCash(amount float64) {
	p.cash = amount
}

func (p Portfolio) GetCash() float64 {
	return p.cash
}
