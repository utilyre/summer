package summer

import (
	"context"
	"crypto/rand"
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestSummer_Sum(t *testing.T) {
	fsys, err := newMockFS()
	if err != nil {
		t.Fatal(err)
	}

	s, err := New(WithFS(fsys), WithRecursive(true))
	if err != nil {
		t.Fatal(err)
	}

	checksums, err := s.Sum(context.Background(), ".")
	if err != nil {
		t.Fatal(err)
	}

	for cs := range checksums {
		if cs.Err != nil {
			t.Error(err)
			continue
		}
	}
}

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

	b.ResetTimer()

	for range b.N {
		checksums, err := s.Sum(ctx, ".")
		if err != nil {
			b.Error(err)
		}

		var checksum Checksum
		for cs := range checksums {
			checksum = cs
		}
		globalChecksum = checksum // to avoid compiler optimization
	}
}

var globalChecksum Checksum

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
