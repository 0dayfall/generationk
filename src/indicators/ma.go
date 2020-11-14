package indicators

type Average struct {
	Comparison
}

func (m *Average) sma(period int) func(float64) float64 {
	var i int
	var sum float64
	var storage = make([]float64, 0, period)

	return func(input float64) (avrg float64) {
		if len(storage) < period {
			sum += input
			storage = append(storage, input)
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
	var result = make([]float64, 0, len(series))
	for _, x := range series {
		result = append(result, ma(x))
	}
	return result
}
