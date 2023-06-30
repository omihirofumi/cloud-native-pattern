package chap04

import (
	"fmt"
	"testing"
	"time"
)

func TestFunnel(t *testing.T) {
	t.Parallel()
	sources := make([]<-chan int, 0)

	for i := 0; i < 3; i++ {
		ch := make(chan int)
		sources = append(sources, ch)

		go func() {
			defer close(ch)

			for i := 1; i <= 5; i++ {
				ch <- i
				time.Sleep(time.Second)
			}
		}()
	}

	dest := Funnel(sources...)
	for d := range dest {
		fmt.Println(d)
	}
}
