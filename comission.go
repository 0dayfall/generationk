package generationk

//Comission interface can be used to implment any comission scheme
type Comission interface {
	GetComisson(price float64, qty int) float64
}

//FixedComission is an example of a fixded comission with 2 limits
type FixedComission struct {
	fixedComissionBelow500 float64
	fixedComissionOver500  float64
}

//GetComission return the amount for more than or less than 500 pieces of asset
func (f FixedComission) GetComisson(price float64, qty int) float64 {
	if qty < 500 {
		return float64(qty) * f.fixedComissionBelow500
	}

	return float64(qty) * f.fixedComissionOver500
}

//PercentageComission is an example of a comission scheme with a fixed percentage
type PercentageComission struct {
	lowComission    float64
	mediumComission float64
	highComission   float64
}

//Returns the comission based on a percentage of the amount, qty is not used in
//this scheme
func (pc PercentageComission) GetComisson(price float64, qty int) float64 {
	amount := price * float64(qty)

	if amount < 10000 {
		return pc.lowComission * amount
	}

	if amount < 20000 {
		return pc.mediumComission * amount
	}

	return pc.highComission * amount
}
