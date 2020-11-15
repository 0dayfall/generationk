package indicators

// OBV means On Balance Volume.
func OBV(priceData, volumeData []float64) []float64 {

	obv := []float64{volumeData[0]}

	for i, vol := range volumeData[1:] {
		if priceData[i] > priceData[i-1] {
			obv = append(obv, obv[i-1]+vol)
		} else if priceData[i] < priceData[i-1] {
			obv = append(obv, obv[i-1]-vol)
		} else {
			obv = append(obv, obv[i-1])
		}
	}

	return obv
}
