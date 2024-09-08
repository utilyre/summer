package pipeline

import "context"

// A Stage is the smallest unit capable of processing data concurrently.
type Stage[T any] interface {
	Pipe(ctx context.Context, in <-chan T) (out <-chan T)
}

// A StageFunc is an adapter to allow the use of ordinary functions as Stage.
type StageFunc[T any] func(ctx context.Context, in <-chan T) <-chan T

func (f StageFunc[T]) Pipe(ctx context.Context, in <-chan T) <-chan T {
	return f(ctx, in)
}

// A Pipeline allows the output of one Stage to be used as the input of another.
type Pipeline[T any] struct {
	// Cap determines maximum queue size for aggregating channels.
	Cap int

	stages []stageInfo[T]
}

type stageInfo[T any] struct {
	n     int
	stage Stage[T]
}

// Append appends n parallel jobs of stage to pl.
func (pl *Pipeline[T]) Append(n int, stage Stage[T]) {
	if n <= 0 {
		panic("pipeline error: non-positive number of stage jobs")
	}

	pl.stages = append(pl.stages, stageInfo[T]{n, stage})
}

// AppendFunc appends n parallel jobs of stage to pl.
func (pl *Pipeline[T]) AppendFunc(n int, stage StageFunc[T]) {
	pl.Append(n, stage)
}

// Clear removes all stages from pl.
func (pl *Pipeline[T]) Clear() {
	pl.stages = nil
}

func (pl *Pipeline[T]) Pipe(ctx context.Context, in <-chan T) <-chan T {
	for _, info := range pl.stages {
		outs := make([]<-chan T, info.n)
		for i := range info.n {
			outs[i] = info.stage.Pipe(ctx, in)
		}

		in = Aggregate(pl.Cap, outs)
	}

	return in
}
