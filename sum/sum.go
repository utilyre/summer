package sum

import (
	"context"
	"crypto/md5"

	"gihtub.com/utilyre/summer/pkg/pipeline"
)

const (
	numReaders   int = 2
	numDigesters int = 5
)

type Sum [md5.Size]byte

func MD5All(root string) (map[string]Sum, error) {
	pl, err := pipeline.New(
		pipeline.WithPipeFunc(read, numReaders),
		pipeline.WithPipeFunc(digest, numDigesters),
	)
	if err != nil {
		return nil, err
	}

	out := pl.Pipe(context.TODO(), walk(context.TODO(), root))
	m := make(map[string]Sum)
	for v := range out {
		cs := v.(checksum)
		m[cs.path] = cs.sum
	}

	return m, nil

	/* g, ctx := errgroup.WithContext(context.Background())

	pathc := walk(ctx, g, root)

	filecs := make([]<-chan file, 0, numReaders)
	for i := 0; i < numReaders; i++ {
		c := read(ctx, g, pathc)
		filecs = append(filecs, c)
	}
	filec := channel.Merge(filecs...)

	checksumcs := make([]<-chan checksum, 0, numDigesters)
	for i := 0; i < numDigesters; i++ {
		c := digest(ctx, g, filec)
		checksumcs = append(checksumcs, c)
	}
	checksumc := channel.Merge(checksumcs...)

	m := make(map[string]Sum)
	for checksum := range checksumc {
		m[checksum.path] = checksum.sum
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("sum: %w", err)
	}
	return m, nil */
}
