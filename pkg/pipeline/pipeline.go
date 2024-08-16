package pipeline

import "context"

type Pipe interface {
	Pipe(ctx context.Context, in <-chan any) (out <-chan any)
}

type PipeFunc func(ctx context.Context, in <-chan any) <-chan any

func (f PipeFunc) Pipe(ctx context.Context, in <-chan any) <-chan any {
	return f(ctx, in)
}

type Pipeline struct {
	pipes []Pipe
}

func New(pipes ...Pipe) *Pipeline {
	return &Pipeline{pipes: pipes}
}

func (pl *Pipeline) Pipe(ctx context.Context, in <-chan any) <-chan any {
	ch := in

	for _, pipe := range pl.pipes {
		ch = pipe.Pipe(ctx, ch)
	}

	return ch
}
