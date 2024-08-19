package pipeline

import (
	"context"
)

type Pipe interface {
	Pipe(ctx context.Context, in <-chan any) (out <-chan any)
}

type PipeFunc func(ctx context.Context, in <-chan any) <-chan any

func (f PipeFunc) Pipe(ctx context.Context, in <-chan any) <-chan any {
	return f(ctx, in)
}

type Pipeline struct {
	AggregBufCap int

	pipes []pipeInfo
}

type pipeInfo struct {
	workers int
	pipe    Pipe
}

func (pl *Pipeline) Append(workers int, pipe Pipe) {
	if workers <= 0 {
		panic("pipeline error: non-positive pipe workers")
	}

	pl.pipes = append(pl.pipes, pipeInfo{workers, pipe})
}

func (pl *Pipeline) AppendFunc(workers int, pipe PipeFunc) {
	pl.Append(workers, pipe)
}

func (pl *Pipeline) Clear() {
	pl.pipes = nil
}

func (pl *Pipeline) Pipe(ctx context.Context, in <-chan any) <-chan any {
	for _, info := range pl.pipes {
		outs := make([]<-chan any, info.workers)
		for i := range info.workers {
			outs[i] = info.pipe.Pipe(ctx, in)
		}

		in = Aggregate(pl.AggregBufCap, outs)
	}

	return in
}
