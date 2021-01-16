package generationk

type Params struct {
	Value, Low, High float64
}

func (p *Params) NewParams(value, low, high float64) *Params {
	return &Params{
		Low:  low,
		High: high,
	}
}

func (p Params) GetValue() float64 {
	return p.Value
}
