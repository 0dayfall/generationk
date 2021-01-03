package generationk

import (
	"time"
)

type OhlcConst int

const (
	//Open value will be used for the data
	Open = iota
	//High value will be used for the data
	High
	//Low value will be used for the data
	Low
	//Close value will be used for the data
	Close
	//Volume value will be used for the data
	Volume
	//Default is used when the data is something else than ohlcv
	Default
)

//DataUpdate is used to update the data in the assets
type DataUpdate interface {
	Update(ohlc OHLC)
}

// OHLC data type
type OHLC struct {
	Time                           time.Time
	Open, High, Low, Close, Volume float64
}

// Asset data type
type Asset struct {
	name string
	ohlc []OHLC
}

//NewAsset is used to create a new asset-
func NewAsset(name string, ohlc OHLC) *Asset {
	return &Asset{
		name: name,
		ohlc: []OHLC{ohlc},
	}
}

//Prepend is used to add the newest data first
func prepend(x []OHLC, y OHLC) []OHLC {
	return append([]OHLC{y}, x...)
}

//Update implements interface to be able to get updated by new data
func (a *Asset) Update(ohlc OHLC, size int) {
	//Dont copy more data than we need
	if len(a.ohlc) < size {
		a.ohlc = prepend(a.ohlc, ohlc)
	} else {
		a.ohlc = prepend(a.ohlc[:size], ohlc)
	}
}

func (a *Asset) Current(ohlcValue OhlcConst) float64 {
	var current float64
	switch ohlcValue {
	case Open:
		current = a.ohlc[0].Open
	case High:
		current = a.ohlc[0].High
	case Low:
		current = a.ohlc[0].Low
	case Volume:
		current = a.ohlc[0].Volume
	case Default:
		break
	}

	return current
}

func (a *Asset) Historic(ohlcValue OhlcConst, barCount int) []float64 {
	return a.sliceOf(ohlcValue, barCount)
}

//CloseArray is used to get the close series
func (a *Asset) sliceOf(ohlcValue OhlcConst, barCount int) []float64 {
	s := make([]float64, barCount)
	barCount = min(len(a.ohlc), barCount)

	switch ohlcValue {
	case Open:
		for i := 0; i < barCount; i++ {
			s[i] = a.ohlc[i].Open
		}
	case High:
		for i := 0; i < barCount; i++ {
			s[i] = a.ohlc[i].High
		}
	case Low:
		for i := 0; i < barCount; i++ {
			s[i] = a.ohlc[i].Low
		}
	case Close:
		for i := 0; i < barCount; i++ {
			s[i] = a.ohlc[i].Close
		}
	case Volume:
		for i := 0; i < barCount; i++ {
			s[i] = a.ohlc[i].Volume
		}
	case Default:
		break

	}

	return s
}

//CloseArray is used to get the close series
func (a *Asset) CloseArray() []float64 {
	s := make([]float64, len(a.ohlc))

	if a.ohlc == nil {
		return nil
	}

	for i, ohlc := range a.ohlc {
		s[i] = ohlc.Close
	}

	return s
}

//Close is used to get the close value
func (a *Asset) Open() float64 {
	if a != nil {
		return a.ohlc[0].Open
	}

	return 0.0
}

//Close is used to get the close value
func (a *Asset) High() float64 {
	if a != nil {
		return a.ohlc[0].High
	}

	return 0.0
}

//Close is used to get the close value
func (a *Asset) Low() float64 {
	if a != nil {
		return a.ohlc[0].Low
	}

	return 0.0
}

//Close is used to get the close value
func (a *Asset) Close() float64 {
	if a != nil {
		return a.ohlc[0].Close
	}

	return 0.0
}

//Close is used to get the close value
func (a *Asset) Volume() float64 {
	if a != nil {
		return a.ohlc[0].Volume
	}

	return 0.0
}

func (a Asset) getName() string {
	return a.name
}

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	// fmt.Printf("date1 %v", date1)
	y2, m2, d2 := date2.Date()
	// fmt.Printf("date2 %v", date2)
	return y1 == y2 && m1 == m2 && d1 == d2
}
