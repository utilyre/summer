package sum

import (
	"context"
	"crypto/md5"
	"fmt"

	"gihtub.com/utilyre/summer/channel"
	"golang.org/x/sync/errgroup"
)

const (
	numReaders   int = 2
	numDigesters int = 5
)

type Sum [md5.Size]byte

func MD5All(root string) (map[string]Sum, error) {
	g, ctx := errgroup.WithContext(context.Background())

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
	return m, nil
}
