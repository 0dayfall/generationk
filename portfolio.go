package generationk

import (
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Direction int

//negativeBalanceErr is used to flag a negative carry
var negativeBalanceErr = errors.New("Balance < 0")

//Direction of a trade: long or short
const (
	Long Direction = iota
	Short
)

//The portfolio holds assets: holdings, the portfolio holds a mutext to
//be able to use the same portfolio when testing many assets in parallell but
//updating the account on a single portfolio
type Portfolio struct {
	m        sync.Mutex
	holdings []Holding
	cash     float64
}

//This is what we are owning, a holding
type Holding struct {
	qty       int
	assetName string
	price     float64
	time      time.Time
}

//Is used to create a new portfolio
func NewPortfolio() *Portfolio {
	portfolio := Portfolio{
		holdings: make([]Holding, 0),
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

//Remove a holding, its sold
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

//AddHolding, its been bought
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

//checkBalance is used to check the balance before buying
func (p *Portfolio) checkBalance(cost float64) error {
	balance := p.cash + cost
	if balance < 0 {
		return negativeBalanceErr
	}
	return nil
}

//addToBalance is used to add to the account after selling with profit
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

//SetBalance is used to set the starting balance of the account
func (p *Portfolio) SetBalance(amount float64) {
	p.m.Lock()
	defer p.m.Unlock()
	p.cash = amount
}

//GetBalance returns in the balance of the account
func (p *Portfolio) GetBalance() float64 {
	return p.cash
}
