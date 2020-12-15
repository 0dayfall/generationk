// +build test

package internal

import (
	genk "generationk/internal"
	"testing"
	"time"
)

func TestAddAsset(t *testing.T) {
	var c genk.Context
	c.AddAsset(genk.NewAsset("Test", nil))

	v := len(c.Asset)
	want := 1

	if got := v; got != want {
		t.Errorf("ReadCSVFile(\"ABB.csv\") = %d, want %d", got, want)
	}
}

func TestIncTime(t *testing.T) {
	c := genk.NewContext()
	startTime := time.Now()
	c.AddStartDate(startTime)
	v := c.Time()
	want := startTime.AddDate(0, 0, 1)

	if got := v; got != want {
		t.Errorf("The days does not match")
	}
}
