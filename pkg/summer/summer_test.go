package summer_test

import (
	"context"
	"crypto/rand"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/utilyre/summer/pkg/summer"
)

func TestSummer_Sum(t *testing.T) {
	s, err := summer.New(
		summer.WithFS(newMockFS(t)),
		summer.WithRecursive(true),
	)
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
	s, err := summer.New(
		summer.WithFS(newMockFS(b)),
		summer.WithRecursive(true),
	)
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

		var checksum summer.Checksum
		for cs := range checksums {
			checksum = cs
		}
		globalChecksum = checksum // to avoid compiler optimization
	}
}

var globalChecksum summer.Checksum

func newMockFS(tb testing.TB) fs.FS {
	foo := make([]byte, 32*1024*1024)
	if _, err := rand.Read(foo); err != nil {
		tb.Fatal("newMockFS:", err)
	}

	bar := make([]byte, 32*1024*1024)
	if _, err := rand.Read(bar); err != nil {
		tb.Fatal("newMockFS:", err)
	}

	baz := make([]byte, 32*1024*1024)
	if _, err := rand.Read(baz); err != nil {
		tb.Fatal("newMockFS:", err)
	}

	return fstest.MapFS{
		"foo": {Data: foo},
		"bar": {Data: bar},
		"baz": {Data: baz},
	}
}
