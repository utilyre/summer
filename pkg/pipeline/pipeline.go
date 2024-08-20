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
	n    int
	pipe Pipe
}

// Append appends n parallel jobs of pipe to pl.
func (pl *Pipeline) Append(n int, pipe Pipe) {
	if n <= 0 {
		panic("pipeline error: non-positive pipe jobs")
	}

	pl.pipes = append(pl.pipes, pipeInfo{n, pipe})
}

// AppendFunc appends n parallel jobs of pipe to pl.
func (pl *Pipeline) AppendFunc(n int, pipe PipeFunc) {
	pl.Append(n, pipe)
}

// Clear removes all pipes from pl.
func (pl *Pipeline) Clear() {
	pl.pipes = nil
}

func (pl *Pipeline) Pipe(ctx context.Context, in <-chan any) <-chan any {
	for _, info := range pl.pipes {
		outs := make([]<-chan any, info.n)
		for i := range info.n {
			outs[i] = info.pipe.Pipe(ctx, in)
		}

		in = Aggregate(pl.AggregBufCap, outs)
	}

	return in
}
