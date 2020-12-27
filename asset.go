package generationk

import (
	"time"

	log "github.com/sirupsen/logrus"
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

// Asset data type
type Asset struct {
	name string
	ohlc []OHLC
}

/*func (a Asset) length() int {
	return len(a.ohlc)
}*/

//NewAsset is used to create a new asset-
func NewAsset(name string, ohlc OHLC) *Asset {
	return &Asset{
		name: name,
		ohlc: []OHLC{ohlc}}
}

//DataUpdate is used to update the data in the assets
type DataUpdate interface {
	Update(ohlc OHLC)
}

// OHLC data type
type OHLC struct {
	time                   time.Time
	open, high, low, close float64
	volume                 int
}

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	// fmt.Printf("date1 %v", date1)
	y2, m2, d2 := date2.Date()
	// fmt.Printf("date2 %v", date2)
	return y1 == y2 && m1 == m2 && d1 == d2
}

func prepend(x []OHLC, y OHLC) []OHLC {
	return append([]OHLC{y}, x...)
}

/*func resize(z []OHLC, period int) []OHLC {
	return z[:len(z)-period]
}*/

//Update interface to be able to get updated by the event queue
func (a *Asset) Update(ohlc OHLC) {
	a.ohlc = prepend(a.ohlc, ohlc)
	/*log.WithFields(log.Fields{
		"a.Ohlc": a.Ohlc,
	}).Debug("ASSET> UPDATE")*/
	/*log.WithFields(log.Fields{
		"a.Ohlc": a.Ohlc,
	}).Debug("ASSET>PREPEND> UPDATE")*/
}

func (a *Asset) GetData(ohlcValue OhlcConst, period int) []float64 {
	switch ohlcValue {
	case Open:
		fallthrough
	case High:
		fallthrough
	case Low:
		fallthrough
	case Default:
		fallthrough
	case Close:
		return a.sliceOfCloseArray(period)
	}
	return nil
}

/*func (a *Asset) Shift(time time.Time) (int, error) {
	var i int
	for ok := true; ok; ok = a.Ohlc[0].Time.Before(time) && len(a.Ohlc) > 0 {
		a.Ohlc = a.Ohlc[1:]
		i++
	}
	return i, nil
}*/

//CloseArray is used to get the close series
func (a *Asset) sliceOfCloseArray(period int) []float64 {
	s := make([]float64, period)

	if a.ohlc == nil {
		return nil
	}

	period = Min(len(a.ohlc), period)
	for i := 0; i < period; i++ {
		s[i] = a.ohlc[i].close
	}
	/*log.WithFields(log.Fields{
		"CloseArray() length": len(s),
	}).Debug("ASSET> CloseArray()")*/
	return s
}

//CloseArray is used to get the close series
func (a *Asset) CloseArray() []float64 {
	s := make([]float64, len(a.ohlc))

	if a.ohlc == nil {
		return nil
	}

	for i, ohlc := range a.ohlc {
		s[i] = ohlc.close
	}
	/*log.WithFields(log.Fields{
		"CloseArray() length": len(s),
	}).Debug("ASSET> CloseArray()")*/
	return s
}

//Close is used to get the close value
func (a *Asset) Close() float64 {
	if a != nil {
		return a.ohlc[0].close
	}
	log.Info("ASSET = NIL")
	return 0.0
}

//CloseAtBar is used to get the close value
func (a *Asset) CloseAtBar(ix int) float64 {
	return a.ohlc[ix].close
}
