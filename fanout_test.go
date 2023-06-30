package chap04

import (
	"fmt"
	"sync"
	"testing"
)

func TestSplit(t *testing.T) {
	t.Parallel()
	source := make(chan int)
	dests := Split(source, 5)

	go func() {
		for i := 1; i <= 10; i++ {
			source <- i
		}

		close(source)
	}()

	var wg sync.WaitGroup
	wg.Add(len(dests))

	for i, ch := range dests {
		go func(i int, d <-chan int) {
			defer wg.Done()

			for val := range d {
				fmt.Printf("#%d got %d\n", i, val)
			}
		}(i, ch)
	}
	wg.Wait()
}
