package pipeline

import (
	"maps"
	"slices"
	"testing"
)

func TestAggregate(t *testing.T) {
	tests := map[string][][]int{
		"no chan": {},
		"one chan": {
			{1, 2, 3, 1, 5}, // chan #1
		},
		"multi chan": {
			{1, 2, 3, 1, 5}, // chan #1
			{8, 2, 2, 3, 4}, // chan #2
			{3, 1, 2, 5, 1}, // chan #3
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var chans []<-chan int
			for _, arr := range test {
				ch := make(chan int)
				chans = append(chans, ch)

				go func() {
					defer close(ch)
					for _, x := range arr {
						ch <- x
					}
				}()
			}

			want := make(map[int]int)
			for _, arr := range test {
				for _, x := range arr {
					want[x]++
				}
			}

			for x := range Aggregate(0, chans) {
				if _, exists := want[x]; !exists {
					t.Error("received unsupplied value:", x)
					continue
				}

				want[x]--
				if want[x] <= 0 {
					delete(want, x)
				}
			}

			if len(want) > 0 {
				s := slices.Collect(maps.Keys(want))
				t.Fatal("missed following values:", s)
			}
		})
	}
}
