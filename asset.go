package generationk

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Asset data type
type Asset struct {
	Name string
	Ohlc []OHLC
}

func (a Asset) length() int {
	return len(a.Ohlc)
}

//NewAsset is used to create a new asset-
func NewAsset(name string, ohlc OHLC) *Asset {
	return &Asset{
		Name: name,
		Ohlc: []OHLC{ohlc}}
}

//DataUpdate is used to update the data in the assets
type DataUpdate interface {
	Update(ohlc OHLC)
}

// OHLC data type
type OHLC struct {
	Time                   time.Time
	Open, High, Low, Close float64
	Volume                 int
}

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	//fmt.Printf("date1 %v", date1)
	y2, m2, d2 := date2.Date()
	//fmt.Printf("date2 %v", date2)
	return y1 == y2 && m1 == m2 && d1 == d2
}

func prepend(x []OHLC, y OHLC) []OHLC {
	return append([]OHLC{y}, x...)
}

func resize(z []OHLC, period int) []OHLC {
	return z[:len(z)-period]
}

//Update interface to be able to get updated by the event queue
func (a *Asset) Update(ohlc OHLC) {
	a.Ohlc = prepend(a.Ohlc, ohlc)
	/*log.WithFields(log.Fields{
		"a.Ohlc": a.Ohlc,
	}).Debug("ASSET> UPDATE")*/
	/*log.WithFields(log.Fields{
		"a.Ohlc": a.Ohlc,
	}).Debug("ASSET>PREPEND> UPDATE")*/

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
func (a *Asset) CloseArray() []float64 {
	s := make([]float64, len(a.Ohlc))

	if a.Ohlc == nil {
		return nil
	}
	for i, ohlc := range a.Ohlc {
		s[i] = ohlc.Close
	}
	/*log.WithFields(log.Fields{
		"CloseArray() length": len(s),
	}).Debug("ASSET> CloseArray()")*/
	return s
}

//Close is used to get the close value
func (a *Asset) Close() float64 {
	if a != nil {
		return a.Ohlc[0].Close
	}
	log.Info("ASSET = NIL")
	return 0.0
}

//CloseAtBar is used to get the close value
func (a *Asset) CloseAtBar(ix int) float64 {
	return a.Ohlc[ix].Close
}
