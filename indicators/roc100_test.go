package generationk

import (
	"testing"
)

func TestAbs(t *testing.T) {

	data := []float64{12.0, 16.5, 17.5, 18.5, 19.5, 20.5}
	roc100 := []float64{0, 0, 0, 19.35483870967742, 18.181818181818183, 17.142857142857142}
	output := RateOfChange100(data, 4)

	for i := 4; i < len(data); i++ {
		if roc100[i] != output[i] {
			t.Errorf("roc100 = %f; want %f", output[i], roc100[i])
		}
	}
}

func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//rand.Int()
	}
}
