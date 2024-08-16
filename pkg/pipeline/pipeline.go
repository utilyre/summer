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

type Option func(o *options) error

type options struct {
	pipes []pipeInfo
}

type pipeInfo struct {
	pipe    Pipe
	workers int
}

func WithPipe(pipe Pipe, workers int) Option {
	return func(o *options) error {
		o.pipes = append(o.pipes, pipeInfo{pipe, workers})
		return nil
	}
}

func WithPipeFunc(pipe PipeFunc, workers int) Option {
	return WithPipe(pipe, workers)
}

func New(opts ...Option) (*Pipeline, error) {
	o := &options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	return &Pipeline{pipes: o.pipes}, nil
}

func (pl *Pipeline) Pipe(ctx context.Context, in <-chan any) <-chan any {
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
