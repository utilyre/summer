package pipeline

import "context"

// A Pipe is the smallest unit capable of processing data concurrently.
type Pipe[T any] interface {
	Pipe(ctx context.Context, in <-chan T) (out <-chan T)
}

// A PipeFunc is an adapter to allow the use of ordinary functions as Pipe.
type PipeFunc[T any] func(ctx context.Context, in <-chan T) <-chan T

func (f PipeFunc[T]) Pipe(ctx context.Context, in <-chan T) <-chan T {
	return f(ctx, in)
}

// A Pipeline allows the output of one Pipe to be used as the input of another.
type Pipeline[T any] struct {
	// Cap determines maximum queue size for each pipe.
	Cap int

	pipes []pipeInfo[T]
}

type pipeInfo[T any] struct {
	n    int
	pipe Pipe[T]
}

// Append appends n parallel jobs of pipe to pl.
func (pl *Pipeline[T]) Append(n int, pipe Pipe[T]) {
	if n <= 0 {
		panic("pipeline error: non-positive pipe jobs")
	}

	pl.pipes = append(pl.pipes, pipeInfo[T]{n, pipe})
}

// AppendFunc appends n parallel jobs of pipe to pl.
func (pl *Pipeline[T]) AppendFunc(n int, pipe PipeFunc[T]) {
	pl.Append(n, pipe)
}

// Clear removes all pipes from pl.
func (pl *Pipeline[T]) Clear() {
	pl.pipes = nil
}

func (pl *Pipeline[T]) Pipe(ctx context.Context, in <-chan T) <-chan T {
	for _, info := range pl.pipes {
		outs := make([]<-chan T, info.n)
		for i := range info.n {
			outs[i] = info.pipe.Pipe(ctx, in)
		}

		in = Aggregate(pl.Cap, outs)
	}

	return in
}
