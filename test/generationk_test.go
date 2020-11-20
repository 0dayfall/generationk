package generationk

import (
	"testing"
	internal "generationk/internal"
	genk "generationk"
	"time"
)

func TestRun(t *testing.T) {
	var c internal.Context
	dm := internal.NewDataManager()
	
	abb := dm.ReadCSVFile("Queue_test.csv")
	c.AddAsset(&abb)
	var strat internal.Strategy
	c.AddStrategy(&strat)
	c.AddStartDate(time.Now())
	c.AddEndDate(time.Now())
	eventChannel := make(chan genk.Event,2)

	genk.Run(&c, eventChannel)
	genk.Run(&c, eventChannel)
}