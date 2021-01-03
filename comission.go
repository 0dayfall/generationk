package generationk

//Comission interface can be used to implment any comission scheme
type Comission interface {
	GetComisson(amount float64, qty int) float64
}

//FixedComission is an example of a fixded comission with 2 limits
type FixedComission struct {
	fixedComissionBelow500 float64
	fixedComissionOver500  float64
}

//GetComission return the amount for more than or less than 500 pieces of asset
func (f FixedComission) GetComisson(amount float64, qty int) float64 {
	if qty < 500 {
		return float64(qty) * f.fixedComissionBelow500
	}

	return float64(qty) * f.fixedComissionOver500
}

//PercentageComission is an example of a comission scheme with a fixed percentage
type PercentageComission struct {
	comission float64
}

//Returns the comission based on a percentage of the amount, qty is not used in
//this scheme
func (pc PercentageComission) GetComisson(amount float64, qty int) float64 {
	return amount * pc.comission
}
