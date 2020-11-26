package indicators

type DEFAULT int

type Average struct {
	*IndicatorStruct
}

// Sma calculates simple moving average of a slice for a certain
// number of time periods.
func (slice mfloat) SMA(period int) []float64 {

	var smaSlice []float64

	for i := period; i <= len(slice); i++ {
		smaSlice = append(smaSlice, Sum(slice[i-period:i])/float64(period))
	}

	return smaSlice
}

// Ema calculates exponential moving average of a slice for a certain
// number of tiSmame periods.
func (slice mfloat) EMA(period int) []float64 {

	var emaSlice []float64

	ak := period + 1
	k := float64(2) / float64(ak)

	emaSlice = append(emaSlice, slice[0])

	for i := 1; i < len(slice); i++ {
		emaSlice = append(emaSlice, (slice[i]*float64(k))+(emaSlice[i-1]*float64(1-k)))
	}

	return emaSlice
}

//SimpleMovingAverage bla bla
func SimpleMovingAverage(series []float64, period int) *Average {
	ma := &Average{
		IndicatorStruct: &IndicatorStruct{},
	}
	ma.Sma(series, period)
	return ma
}

func (m *Average) sma(period int) func(float64) float64 {
	var i int
	var sum float64
	var storage = make([]float64, period)

	return func(input float64) (avrg float64) {
		if len(storage) < period {
			sum += input
			storage[i] = input
		}

		sum += input - storage[i]
		storage[i], i = input, (i+1)%period
		avrg = sum / float64(len(storage))

		return
	}
}

//Sma function is used to calc moving averages
func (m *Average) Sma(series []float64, period int) []float64 {
	ma := m.sma(period)
	var result = make([]float64, len(series))
	for i, x := range series {
		result[i] = ma(x) //append(result, ma(x))
	}
	m.IndicatorStruct.defaultValues = result

	return result
}
