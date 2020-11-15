package indicators

// Ichimoku Cloud.
func IchimokuCloud(priceData, lowData, highData mfloat, configs []int) ([]float64, []float64, []float64, []float64, []float64) {

	var conversionLine, baseLine, leadSpanA, leadSpanB, lagSpan []float64

	conversionLine = DivSlice(SubSlices(highData.SMA(9), lowData.SMA(9)), 2)
	baseLine = DivSlice(SubSlices(highData.SMA(26), lowData.SMA(26)), 2)
	leadSpanA = DivSlice(AddSlices(conversionLine, baseLine), 2)
	leadSpanB = DivSlice(SubSlices(highData.SMA(52), lowData.SMA(52)), 2)
	lagSpan = priceData[0 : len(priceData)-26]

	return conversionLine, baseLine, leadSpanA, leadSpanB, lagSpan
}
