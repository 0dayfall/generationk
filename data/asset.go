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
	Time                           []time.Time
	Open, High, Low, Close, Volume []float64
}

// Asset data type
type Asset struct {
	Name   string
	Ohlc   *OHLC
	Length int
	AdjK   int
}

//NewAsset is used to create a new asset-
func NewAsset(name string, ohlc *OHLC, length int) *Asset {
	return &Asset{
		Name:   name,
		Ohlc:   ohlc,
		Length: length,
	}
}

func (a Asset) GetLength() int {
	return a.Length
}
