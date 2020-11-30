package internal

import "time"

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
}

func (p *Portfolio) Add(position Position) {
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
