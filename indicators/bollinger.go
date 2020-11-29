package indicators

// BollingerBands returns upper band, lower band and simple moving
// average of a slice.
func BollingerBands(slice []float64, period int, nStd float64) ([]float64, []float64, []float64) {

	var upperBand, lowerBand, middleBand []float64

	middleBand = SimpleMa(slice, period)
	std := Std(middleBand)
	upperBand = AddToAll(middleBand, std*nStd)
	lowerBand = AddToAll(middleBand, -1.0*std*nStd)

	return middleBand, upperBand, lowerBand
}
