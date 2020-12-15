package internal

import (
	"container/heap"
	"fmt"
	int "generationk/internal"
	"testing"
	"time"
)

func TestOhlcHeap(t *testing.T) {

	//Oldest value is 3 seconds ago
	timex := time.Now().Add(-3 * time.Second)
	h := &int.OhlcHeap{
		//Second oldest
		int.OHLC{Time: time.Now(), Close: 3.0},
		//Oldest value se above
		int.OHLC{Time: timex, Close: 4.0},
		//A little bit later
		int.OHLC{Time: time.Now(), Close: 2.0},
	}
	heap.Init(h)
	//Nest came last
	heap.Push(h, int.OHLC{Time: time.Now(), Close: 1.0})

	//Newest value should come first, imporant to have newest data first
	for h.Len() > 0 {
		fmt.Printf("%d \n", heap.Pop(h))
	}
}
