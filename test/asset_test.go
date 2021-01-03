// +build asset_test

package generationk_test

import (
	K "generationk"
	"testing"
	"time"
)

func TestCreateAsset(t *testing.T) {
	//fmt.Printf("Testing %s", t.Name())
	asset := K.NewAsset("test", K.OHLC{Time: time.Now(), Open: 1.0, High: 2.0, Low: 3.0, Close: 4.0, Volume: 5})

	if asset != nil && isTest(asset) {
		t.Fatalf(`Asset not of type Asset`)
	}
}

func TestCurrentOpen(t *testing.T) {
	//fmt.Printf("Testing %s", t.Name())
	asset := K.NewAsset("test", K.OHLC{Time: time.Now(), Open: 1.0, High: 2.0, Low: 3.0, Close: 4.0, Volume: 5})

	open := asset.Current(K.Open)

	if open != 1.0 {
		t.Fatalf(`Asset not of type Asset`)
	}
}

func TestHistoryOpen(t *testing.T) {
	//fmt.Printf("Testing %s", t.Name())
	asset := K.NewAsset("test", K.OHLC{Time: time.Now(), Open: 1.0, High: 2.0, Low: 3.0, Close: 4.0, Volume: 5})

	open := asset.History(K.Open, 1)

	if open[0] != 1.0 {
		t.Fatalf(`Asset not of type Asset`)
	}
}

func TestGetHigh(t *testing.T) {
	//fmt.Printf("Testing %s", t.Name())
	asset := K.NewAsset("test", K.OHLC{Time: time.Now(), Open: 1.0, High: 2.0, Low: 3.0, Close: 4.0, Volume: 5})

	high := asset.History(K.High, 1)

	if high[0] != 2.0 {
		t.Fatalf(`Asset not of type Asset`)
	}
}

func TestGetLow(t *testing.T) {
	//fmt.Printf("Testing %s", t.Name())
	asset := K.NewAsset("test", K.OHLC{Time: time.Now(), Open: 1.0, High: 2.0, Low: 3.0, Close: 4.0, Volume: 5})

	low := asset.History(K.Low, 1)

	if low[0] != 3.0 {
		t.Fatalf(`Asset not of type Asset`)
	}
}

func TestGetClose(t *testing.T) {
	//fmt.Printf("Testing %s", t.Name())
	asset := K.NewAsset("test", K.OHLC{Time: time.Now(), Open: 1.5, High: 2.0, Low: 1.0, Close: 4.0, Volume: 5})

	close := asset.History(K.Close, 1)

	if close[0] != 4.0 {
		t.Fatalf(`Asset not of type Asset`)
	}
}

func TestGetTwoClose(t *testing.T) {
	//fmt.Printf("Testing %s", t.Name())
	asset := K.NewAsset("test", K.OHLC{Time: time.Now(), Open: 1.5, High: 2.0, Low: 1.0, Close: 4.0, Volume: 5})
	asset.Update(K.OHLC{Time: time.Now(), Open: 1.5, High: 2.0, Low: 1.0, Close: 4.0, Volume: 5}, 2)

	close := asset.History(K.Close, 2)

	if len(close) != 2 {
		t.Fatalf(`Asset should contain array with length 2`)
	}
}

func TestGetVolume(t *testing.T) {
	//fmt.Printf("Testing %s", t.Name())
	asset := K.NewAsset("test", K.OHLC{Time: time.Now(), Open: 1.0, High: 2.0, Low: 3.0, Close: 4.0, Volume: 5})

	volume := asset.History(K.Volume, 1)

	if volume[0] != 5 {
		t.Fatalf(`Asset not of type Asset`)
	}
}

func isTest(t interface{}) bool {
	switch t.(type) {
	case K.Asset:
		return true
	default:
		return false
	}
}
