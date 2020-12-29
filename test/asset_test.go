// +build test

package generationk_test

import (
	K "generationk/Asset"
	"testing"
)

func TestCreateAsset(t *testing.T) {
	asset := K.NewAsset("test", K.OHLC{1.0, 2.0, 3.0, 4.0, 5})

	if asset != nil && isTest(asset) {
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
