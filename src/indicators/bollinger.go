package indicators

// BollingerBands returns upper band, lower band and simple moving
// average of a slice.
func BollingerBands(slice mfloat, period int, nStd float64) ([]float64, []float64, []float64) {

	var upperBand, lowerBand, middleBand mfloat

	middleBand = slice.SMA(period)
	std := Std(middleBand)
	upperBand = middleBand.AddToAll(std * nStd)
	lowerBand = middleBand.AddToAll(-1.0 * std * nStd)

	return middleBand, upperBand, lowerBand
}
