package generationk

type Comission interface {
	GetComisson(amount float64, qty int) float64
}

type FixedComission struct {
	fixedComissionBelow500 float64
	fixedComissionOver500  float64
}

func (f FixedComission) GetComisson(amount float64, qty int) float64 {
	if qty < 500 {
		return float64(qty) * f.fixedComissionBelow500
	}

	return float64(qty) * f.fixedComissionOver500
}

type PercentageComission struct {
	comission float64
}

func (pc PercentageComission) GetComisson(amount float64, qty int) float64 {
	return amount * pc.comission
}
