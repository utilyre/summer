package sum

import (
	"context"
	"crypto/md5"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

type file struct {
	path string
	data []byte
}

type checksum struct {
	path string
	sum  Sum
}

func walk(ctx context.Context, g *errgroup.Group, root string) <-chan string {
	out := make(chan string)

	g.Go(func() error {
		defer close(out)

		return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() {
				return nil
			}

			select {
			case out <- path:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	})

	return out
}

func read(ctx context.Context, g *errgroup.Group, in <-chan string) <-chan file {
	out := make(chan file)

	g.Go(func() error {
		defer close(out)

		for p := range in {
			data, err := os.ReadFile(p)
			if err != nil {
				return err
			}

			select {
			case out <- file{path: p, data: data}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		return nil
	})

	return out
}

func digest(ctx context.Context, g *errgroup.Group, in <-chan file) <-chan checksum {
	out := make(chan checksum)

	g.Go(func() error {
		defer close(out)

		for f := range in {
			select {
			case out <- checksum{path: f.path, sum: md5.Sum(f.data)}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		return nil
	})

	return out
}
