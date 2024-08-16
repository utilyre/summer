package pipeline

import "context"

type Pipe interface {
	Pipe(ctx context.Context, in <-chan any) (out <-chan any)
}

type PipeFunc func(ctx context.Context, in <-chan any) <-chan any

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

/* type Pipe interface {
	Pipe(ctx context.Context, in <-chan any) (out <-chan any)
}

type PipeFunc func(ctx context.Context, in <-chan any) <-chan any

func (f PipeFunc) Pipe(ctx context.Context, in <-chan any) <-chan any {
	return f(ctx, in)
}

type Pipeline struct {
	in    <-chan any
	out   <-chan any
	pipes []pipeInfo
}

type pipeInfo struct {
	pipe PipeFunc
	num  int
}

type Option func(opts *options) error

type options struct {
	pipes []pipeInfo
}

func WithPipe(pipe PipeFunc, num int) Option {
	return func(opts *options) error {
		opts.pipes = append(opts.pipes, pipeInfo{pipe, num})
		return nil
	}
}

func New(in <-chan any, opts ...Option) (*Pipeline, error) {
	o := &options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	return &Pipeline{
		in:    in,
		out:   make(<-chan any),
		pipes: o.pipes,
	}, nil
}

func (pl *Pipeline) Launch(ctx context.Context) {
	ch := pl.in

	for _, info := range pl.pipes {
		outs := make([]<-chan any, info.num)
		for i := 0; i < info.num; i++ {
			outs[i] = info.pipe(ctx, ch)
		}

		ch = channel.Merge(outs...)
	}

	pl.out = ch
} */
