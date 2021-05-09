package generationk

import (
	"testing"
)

func TestLength(t *testing.T) {

	high := []float64{12.0, 16.5, 17.5, 18.5, 17.5, 20.5, 12, 12, 12, 12}
	low := []float64{12.0, 16.5, 17.5, 18.5, 17.5, 20.5, 12, 12, 12, 12}
	up, _, _ := Donchian(high, low, 3)

	if len(up) != len(high) {
		t.Errorf("up length = %d; want %d", len(up), len(high))
	}
}

func TestDonch(t *testing.T) {

	//High
	high := []float64{849.50, 871.00, 897.00, 905.50, 912.00, 946.00, 935.00, 977.50}
	//Low
	low := []float64{820.00, 847.50, 873.00, 883.50, 892.50, 913.50, 906.5, 911.0}

	upper := []float64{849.5, 871, 897, 905.5, 912, 946, 946, 977.5}
	lower := []float64{820.00, 820.00, 820.00, 820.00, 847.5, 873.0, 883.5, 892.5}

	up, _, dn := Donchian(high, low, 4)

	for i := 0; i < len(up); i++ {
		if up[i] != upper[i] {
			t.Errorf("donch up = %f; want %f", up[i], upper[i])
		}
	}

	for i := 0; i < len(dn); i++ {
		if dn[i] != lower[i] {
			t.Errorf("donch dn = %f; want %f", dn[i], lower[i])
		}
	}
}
