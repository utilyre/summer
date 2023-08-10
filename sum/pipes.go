package sum

import (
	"context"
	"crypto/md5"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type result[T any] struct {
	data T
	err  error
}

type file struct {
	path    string
	content []byte
}

type checksum struct {
	path string
	sum  Sum
}

func walk(ctx context.Context, root string) (<-chan string, chan error) {
	out := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(out)

		errc <- filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() {
				return nil
			}

			select {
			case out <- path:
			case <-ctx.Done():
				return errors.New("walk cancelled")
			}
			return nil
		})
	}()

	return out, errc
}

func read(ctx context.Context, paths <-chan string) <-chan result[file] {
	out := make(chan result[file])

	go func() {
		defer close(out)

		for path := range paths {
			rout := result[file]{}
			content, err := os.ReadFile(path)

			if err == nil {
				rout.data = file{
					path:    path,
					content: content,
				}
			} else {
				rout.err = err
			}

			select {
			case out <- rout:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func digest(ctx context.Context, contents <-chan result[file]) <-chan result[checksum] {
	out := make(chan result[checksum])

	go func() {
		defer close(out)

		for rin := range contents {
			rout := result[checksum]{}

			if rin.err == nil {
				rout.data = checksum{
					path: rin.data.path,
					sum:  md5.Sum(rin.data.content),
				}
			} else {
				rout.err = rin.err
			}

			select {
			case out <- rout:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
