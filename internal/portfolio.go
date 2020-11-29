package internal

type Position {
	amount int
	asset *Asset
	time time.Time
	comission float64
}

type Portfolio struct {
	Positions []Position
}

func (p *Portfolio) SetCash(amount float64) {
	b.cash = amount
}

func (p *Portfolio) GetCash() float64{
	return b.cash
}