package internal

import "time"

// Asset data type
type Asset struct {
	Name string
	Ohlc []OHLC
}

//Type is the type for ohlc
type Type struct {
	Open  string
	High  string
	Low   string
	close string
}

// OHLC data type
type OHLC struct {
	Time                   time.Time
	Open, High, Low, Close float64
	Volume                 int
}

//GetCloseArray is used to get the close series
func (a *Asset) GetCloseArray() []float64 {
	s := make([]float64, len(a.Ohlc))

	for _, ohlc := range a.Ohlc {
		s = append(s, ohlc.Close)
	}
	return s
}

//Close is used to get the close value
func (a *Asset) Close() float64 {
	return a.Ohlc[len(a.Ohlc)-1].Close
}

//CloseAtBar is used to get the close value
func (a *Asset) CloseAtBar(ix int) float64 {
	return a.Ohlc[ix].Close
}

// Portfolio structure
type Portfolio struct {
}
