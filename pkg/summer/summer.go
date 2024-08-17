package summer

import (
	"context"

	"gihtub.com/utilyre/summer/pkg/pipeline"
)

type Algorithm int

const (
	AlgorithmMD5 Algorithm = iota + 1
)

func SumTree(
	ctx context.Context,
	root string,
	algo Algorithm,
) ([]ChecksumInfo, error) {
	var pl pipeline.Pipeline
	pl.Append( /* TODO: config */ 2, readerPipe{})
	pl.Append( /* TODO: config */ 5, digesterPipe{})
	out := pl.Pipe(ctx, walkerPipe{root: root}.Pipe(ctx, nil))

	var checksums []ChecksumInfo
	for v := range out {
		info := v.(ChecksumInfo)
		checksums = append(checksums, info)
	}

	return checksums, nil
}
