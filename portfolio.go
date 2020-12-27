package generationk

import (
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Direction int

var negativeBalanceErr = errors.New("Balance < 0")

const (
	Long Direction = iota
	Short
)

type Portfolio struct {
	m        sync.Mutex
	holdings []Holding
	cash     float64
}

type Holding struct {
	qty       int
	assetName string
	price     float64
	time      time.Time
}

func NewPortfolio() *Portfolio {
	portfolio := Portfolio{
		holdings: make([]Holding, 0),
		cash:     0,
	}
	return &portfolio
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

	for k := range p.holdings {
		if p.holdings[k].assetName == assetName {

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
		"asset": position.assetName,
		"time":  position.time,
		"price": position.price,
		"qty":   position.qty,
	}).Info("PORTFOLIO> Removing position from portfolio")

	pos := -1

	for k := range p.holdings {
		if position.assetName == p.holdings[k].assetName {
			pos = k
		}
	}

	p.holdings = remove(pos, p.holdings)
}

func remove(ix int, holdings []Holding) []Holding {
	return append(holdings[:ix], holdings[ix+1:]...)
}

func (p *Portfolio) AddHolding(position Holding) {
	p.m.Lock()
	defer p.m.Unlock()
	log.WithFields(log.Fields{
		"asset": position.assetName,
		"time":  position.time,
		"Qty":   position.qty,
	}).Info("PORTFOLIO> Adding position to portfolio")

	p.holdings = append(p.holdings, position)
}

func (p *Portfolio) checkBalance(cost float64) error {
	balance := p.cash + cost
	if balance < 0 {
		return negativeBalanceErr
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

func (p *Portfolio) SetBalance(amount float64) {
	p.m.Lock()
	defer p.m.Unlock()
	p.cash = amount
}

func (p *Portfolio) GetBalance() float64 {
	return p.cash
}
