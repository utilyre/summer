package summer_test

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"io/fs"
	"math/rand"
	"testing"
	"testing/fstest"

	"github.com/utilyre/summer/pkg/summer"
)

func TestSummer_Sum(t *testing.T) {
	s, err := summer.New(
		summer.WithFS(newMockFS(t, 10)),
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
	benchmarkSummer_Sum(b, false)
}

func BenchmarkSummer_Sum_recursive(b *testing.B) {
	benchmarkSummer_Sum(b, true)
}

func benchmarkSummer_Sum(b *testing.B, recursive bool) {
	ctx := context.Background()
	b.ResetTimer()

	for range b.N {
		b.StopTimer()
		s, err := summer.New(
			summer.WithFS(newMockFS(b, 100)),
			summer.WithRecursive(recursive),
		)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()

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

func newMockFS(tb testing.TB, numFiles int) fs.FS {
	fsys := fstest.MapFS{}

	for i := range numFiles {
		size := (rand.Intn(1<<10-1) + 1) << 10 // [1kB, 1MB)
		data := make([]byte, size)
		if _, err := crand.Read(data); err != nil {
			tb.Fatal("newMockFS:", err)
		}
		fsys[fmt.Sprintf("file_%03d", i)] = &fstest.MapFile{Data: data}
	}

	return fsys
}
