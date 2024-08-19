package pipeline

import "reflect"

func Aggregate[T any](cap int, cs []<-chan T) <-chan T {
	out := make(chan T, cap)

	go func() {
		defer close(out)

		cases := make([]reflect.SelectCase, len(cs))
		for i, c := range cs {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			}
		}

		numClosed := 0
		for numClosed < len(cases) {
			idx, v, open := reflect.Select(cases)
			if !open {
				cases[idx].Chan = reflect.Value{}
				numClosed++
				continue
			}

			out <- v.Interface().(T)
		}
	}()

	return out
}
