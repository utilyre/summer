package sum

import (
	"context"
	"crypto/md5"
	"fmt"

	"gihtub.com/utilyre/summer/channel"
)

const (
	numReaders   int = 5
	numDigesters int = 5
)

type Sum [md5.Size]byte

func MD5All(root string) (map[string]Sum, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errcs := make([]<-chan error, 0, 1+numReaders+numDigesters)

	pathc, errc := walk(ctx, root)
	errcs = append(errcs, errc)

	filecs := make([]<-chan file, 0, numReaders)
	for i := 0; i < numReaders; i++ {
		c, errc := read(ctx, pathc)
		filecs = append(filecs, c)
		errcs = append(errcs, errc)
	}
	filec := channel.Merge(filecs...)

	checksumcs := make([]<-chan checksum, 0, numDigesters)
	for i := 0; i < numDigesters; i++ {
		c, errc := digest(ctx, filec)
		checksumcs = append(checksumcs, c)
		errcs = append(errcs, errc)
	}
	checksumc := channel.Merge(checksumcs...)
	errc = channel.Merge(errcs...)

	m := make(map[string]Sum)
	for checksum := range checksumc {
		m[checksum.path] = checksum.sum
	}

	for err := range errc {
		if err != nil {
			return nil, fmt.Errorf("sum: %w", err)
		}
	}
	return m, nil
}
