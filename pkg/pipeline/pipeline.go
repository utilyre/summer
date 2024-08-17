package pipeline

import (
	"context"
	"sync"
)

type Pipe interface {
	Pipe(ctx context.Context, in <-chan any) (out <-chan any)
}

type PipeFunc func(ctx context.Context, in <-chan any) <-chan any

func (f PipeFunc) Pipe(ctx context.Context, in <-chan any) <-chan any {
	return f(ctx, in)
}

type Pipeline struct {
	pipes []pipeInfo
}

type pipeInfo struct {
	workers int
	pipe    Pipe
}

func (pl *Pipeline) Append(workers int, pipe Pipe) {
	if workers <= 0 {
		// TODO: write test for this
		panic("pipeline error: non-positive pipe workers")
	}

	pl.pipes = append(pl.pipes, pipeInfo{workers, pipe})
}

func (pl *Pipeline) AppendFunc(workers int, pipe PipeFunc) {
	pl.Append(workers, pipe)
}

func (pl Pipeline) Pipe(ctx context.Context, in <-chan any) <-chan any {
	for _, info := range pl.pipes {
		outs := make([]<-chan any, info.workers)
		for i := 0; i < info.workers; i++ {
			outs[i] = info.pipe.Pipe(ctx, in)
		}

		in = aggregateChans(outs)
	}

	return in
}

func aggregateChans(cs []<-chan any) <-chan any {
	out := make(chan any)

	var wg sync.WaitGroup
	wg.Add(len(cs))

	for _, c := range cs {
		go func(c <-chan any) {
			for n := range c {
				out <- n
			}

			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
