package generationk

import (
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Direction int

const (
	Long Direction = iota
	Short
)

type Portfolio struct {
	m        sync.Mutex
	Holdings []Holding
	cash     float64
}

type Holding struct {
	Qty       int
	AssetName string
	Price     float64
	Time      time.Time
}

func NewPortfolio() *Portfolio {
	portfolio := Portfolio{
		Holdings: make([]Holding, 0),
		cash:     0,
	}
	return &portfolio
}

//IsOwning is used to find out if a position is already owned in this asset
func (p *Portfolio) IsOwning(assetName string) bool {
	p.m.Lock()
	defer p.m.Unlock()

	for k := range p.Holdings {
		if p.Holdings[k].AssetName == assetName {
			log.WithFields(log.Fields{
				"AssetName": assetName,
			}).Debug("Already owned")
			return true
		}
	}
	return false
}

func (p *Portfolio) RemoveHolding(position Holding) {
	p.m.Lock()
	defer p.m.Unlock()

	log.WithFields(log.Fields{
		"asset": position.AssetName,
		"time":  position.Time,
		"price": position.Price,
		"qty":   position.Qty,
	}).Info("PORTFOLIO> Removing position from portfolio")

	pos := -1
	for k := range p.Holdings {
		if position.AssetName == p.Holdings[k].AssetName {
			pos = k
		}
	}
	p.Holdings = remove(pos, p.Holdings)
}

func remove(ix int, holdings []Holding) []Holding {
	return append(holdings[:ix], holdings[ix+1:]...)
}

func (p *Portfolio) AddHolding(position Holding) {
	p.m.Lock()
	defer p.m.Unlock()
	log.WithFields(log.Fields{
		"asset": position.AssetName,
		"time":  position.Time,
		"Qty":   position.Qty,
	}).Debug("PORTFOLIO> Adding position to portfolio")
	p.Holdings = append(p.Holdings, position)
}

func (p *Portfolio) checkBalance(cost float64) error {
	balance := p.cash + cost
	if balance < 0 {
		return errors.New("Balance < 0")
	}
	return nil
}

func (p *Portfolio) addToBalance(value float64) {
	defer p.m.Unlock()
	p.m.Lock()
	p.cash += value
}

//SubtractFromBalance is used to decrease the amount on the account
func (p *Portfolio) subtractFromBalance(cost float64) error {
	p.m.Lock()
	defer p.m.Unlock()
	err := p.checkBalance(cost)
	if err != nil {
		return err
	}
	p.cash -= cost
	return nil
}

func (p *Portfolio) SetCash(amount float64) {
	p.cash = amount
}

func (p *Portfolio) GetCash() float64 {
	return p.cash
}
