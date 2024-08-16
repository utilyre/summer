package sum

import (
	"context"
	"crypto/md5"
	"io/fs"
	"os"
	"path/filepath"
)

type file struct {
	path string
	data []byte
}

type checksum struct {
	path string
	sum  Sum
}

func walk(ctx context.Context, root string) <-chan any {
	out := make(chan any)

	go func() {
		defer close(out)

		if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
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
		}); err != nil {
			panic("TODO")
		}
	}()

	return out
}

func read(ctx context.Context, in <-chan any) <-chan any {
	out := make(chan any)

	go func() {
		defer close(out)

		for v := range in {
			name := v.(string)
			data, err := os.ReadFile(name)
			if err != nil {
				panic("TODO")
			}

			select {
			case out <- file{path: name, data: data}:
			case <-ctx.Done():
				panic("TODO")
			}
		}
	}()

	return out
}

func digest(ctx context.Context, in <-chan any) <-chan any {
	out := make(chan any)

	go func() {
		defer close(out)

		for v := range in {
			f := v.(file)
			select {
			case out <- checksum{path: f.path, sum: md5.Sum(f.data)}:
			case <-ctx.Done():
				panic("TODO")
			}
		}
	}()

	return out
}
