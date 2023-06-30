package chap04

func Split(source <-chan int, n int) []<-chan int {
	dests := make([]<-chan int, 0)

	for i := 0; i < n; i++ {
		ch := make(chan int)
		dests = append(dests, ch)

		go func() {
			defer close(ch)

			for val := range source {
				ch <- val
			}
		}()
	}
	return dests
}
