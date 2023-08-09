package channel

import "sync"

func Merge[T any](cs ...<-chan T) <-chan T {
	out := make(chan T)

	var wg sync.WaitGroup
	wg.Add(len(cs))

	capture := func(c <-chan T) {
		for n := range c {
			out <- n
		}

		wg.Done()
	}

	for _, c := range cs {
		go capture(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
