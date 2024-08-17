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
	var pl pipeline.Pipeline
	pl.AppendFunc(numReaders, read)
	pl.AppendFunc(numDigesters, digest)

	out := pl.Pipe(context.TODO(), walk(context.TODO(), root))
	m := make(map[string]Sum)
	for v := range out {
		cs := v.(checksum)
		m[cs.path] = cs.sum
	}

	return m, nil
}
