package generationk

import (
	genk "generationk"
	internal "generationk/internal"
	"sync"
	"testing"
	"time"
)

func TestEventQueue(t *testing.T) {
	var c internal.Context
	dm := internal.NewDataManager()

	abb := dm.ReadCSVFile("Queue_test.csv")
	c.AddAsset(&abb)
	c.AddStrategy(nil)
	c.AddStartDate(time.Now())
	c.AddEndDate(time.Now())
	channelSize := len(abb.CloseArray())

	channel := make(chan internal.OHLC, channelSize)
	go genk.PutData(&c, channel)

	var wg sync.WaitGroup
	var value internal.OHLC
	wg.Add(1)
	go func() {
		value = <-channel
		wg.Done()
	}()

	wg.Wait()
	if value.Close == abb.CloseArray()[0] {
		t.Errorf("The value read in the .csv file is not the same as fixture")
	}

}

func TestOHLCQueue(t *testing.T) {
	var c internal.Context
	dm := internal.NewDataManager()

	abb := dm.ReadCSVFile("Queue_test.csv")
	c.AddAsset(&abb)
	c.AddStrategy(nil)
	c.AddStartDate(time.Now())
	c.AddEndDate(time.Now())
	channelSize := len(abb.CloseArray())

	channel := make(chan internal.OHLC, channelSize)
	go genk.PutData(&c, channel)

	valueComparison := internal.OHLC{Time: time.Now(), Open: 0, High: 0, Low: 0, Close: 229.20, Volume: 0}
	var value internal.OHLC

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		value = <-channel
		wg.Done()
	}()

	wg.Wait()
	if value.Close != valueComparison.Close {
		t.Errorf("The value read in the .csv file is not the same as fixture")
	}

}
