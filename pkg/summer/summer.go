package summer

import (
	"context"

	"gihtub.com/utilyre/summer/pkg/pipeline"
	"golang.org/x/sync/errgroup"
)

type Algorithm int

const (
	AlgorithmMD5 Algorithm = iota + 1
)

func SumTree(
	ctx context.Context,
	root string,
	algo Algorithm,
) ([]Checksum, error) {
	g, ctx := errgroup.WithContext(ctx)

	var pl pipeline.Pipeline
	pl.Append( /* TODO: config */ 2, readerPipe{g})
	pl.Append( /* TODO: config */ 5, digesterPipe{g})
	out := pl.Pipe(ctx, walkerPipe{g, root}.Pipe(ctx, nil))

	var checksums []Checksum
	for v := range out {
		info := v.(Checksum)
		checksums = append(checksums, info)
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return checksums, nil
}
