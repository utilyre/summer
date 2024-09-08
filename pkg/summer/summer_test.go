package summer

import (
	"context"
	"crypto/rand"
	"io/fs"
	"testing"
	"testing/fstest"
)

func BenchmarkSummer_Sum(b *testing.B) {
	fsys, err := newMockFS()
	if err != nil {
		b.Fatal(err)
	}

	s, err := New(WithFS(fsys), WithRecursive(true))
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	for range b.N {
		if _, err := s.Sum(ctx, "."); err != nil {
			b.Error(err)
		}
	}
}

func newMockFS() (fs.FS, error) {
	foo := make([]byte, 32*1024*1024)
	if _, err := rand.Read(foo); err != nil {
		return nil, err
	}

	bar := make([]byte, 32*1024*1024)
	if _, err := rand.Read(bar); err != nil {
		return nil, err
	}

	baz := make([]byte, 32*1024*1024)
	if _, err := rand.Read(baz); err != nil {
		return nil, err
	}

	return fstest.MapFS{
		"foo": {Data: foo},
		"bar": {Data: bar},
		"baz": {Data: baz},
	}, nil
}
