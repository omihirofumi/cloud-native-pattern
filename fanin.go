package chap04

import "sync"

func Funnel(sources ...<-chan int) <-chan int {
	dest := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(sources))

	for _, ch := range sources {
		go func(c <-chan int) {
			defer wg.Done()

			for n := range c {
				dest <- n
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(dest)
	}()
	return dest
}
