package pipeline

import (
	"context"
)

// A Pipe is the smallest unit capable of processing data concurrently.
type Pipe interface {
	Pipe(ctx context.Context, in <-chan any) (out <-chan any)
}

// A PipeFunc is an adapter to allow the use of ordinary functions as Pipe.
type PipeFunc func(ctx context.Context, in <-chan any) <-chan any

func (f PipeFunc) Pipe(ctx context.Context, in <-chan any) <-chan any {
	return f(ctx, in)
}

// A Pipeline allows the output of one Pipe to be used as the input of another.
type Pipeline struct {
	AggregBufCap int

	pipes []pipeInfo
}

type pipeInfo struct {
	workers int
	pipe    Pipe
}

// Append appends workers number of instances of pipe to pl in parallel.
func (pl *Pipeline) Append(workers int, pipe Pipe) {
	if workers <= 0 {
		panic("pipeline error: non-positive pipe workers")
	}

	pl.pipes = append(pl.pipes, pipeInfo{workers, pipe})
}

// AppendFunc appends workers number of instances of pipe to pl in parallel.
func (pl *Pipeline) AppendFunc(workers int, pipe PipeFunc) {
	pl.Append(workers, pipe)
}

// Clear removes all pipes from pl.
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
