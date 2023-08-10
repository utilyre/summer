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

func walk(ctx context.Context, root string) (<-chan string, <-chan error) {
	out := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
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
		if err != nil {
			errc <- err
			return
		}
	}()

	return out, errc
}

func read(ctx context.Context, in <-chan string) (<-chan file, <-chan error) {
	out := make(chan file)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for p := range in {
			data, err := os.ReadFile(p)
			if err != nil {
				errc <- err
				return
			}

			select {
			case out <- file{path: p, data: data}:
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			}
		}
	}()

	return out, errc
}

func digest(ctx context.Context, in <-chan file) (<-chan checksum, <-chan error) {
	out := make(chan checksum)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for f := range in {
			select {
			case out <- checksum{path: f.path, sum: md5.Sum(f.data)}:
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			}
		}
	}()

	return out, errc
}
