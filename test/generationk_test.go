// +build test
package generationk

import (
	genk "generationk"
	internal "generationk/internal"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	ctx := genk.NewContext()
	dm := genk.NewCSVDataManager(ctx)

	abb := dm.ReadCSVFile("Queue_test.csv")
	c.AddAsset(&abb)
	var strat internal.Strategy
	c.AddStrategy(&strat)
	c.AddStartDate(time.Now())
	c.AddEndDate(time.Now())

	genk.RunLive(&c)

}
