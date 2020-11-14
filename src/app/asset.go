package main

import "time"

// Asset data type
type Asset struct {
	name string
	ohlc []OHLC
}

// OHLC data type
type OHLC struct {
	time                   time.Time
	open, high, low, close float64
	volume                 int
}

func (a *Asset) GetCloseArray() []float64 {
	s := make([]float64, len(a.ohlc))

	for _, ohlc := range a.ohlc {
		s = append(s, ohlc.close)
	}
	return s
}

func (a *Asset) GetClose() float64 {
	return a.ohlc[len(a.ohlc)-1].close
}

// Portfolio structure
type Portfolio struct {
}

//Event type
type Event struct{}

//Tick event type
type Tick Event

//Signal event type
type Signal Event

//Order event type
type Order Event

//Fill event type
type Fill Event
