package summer

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

type walkerPipe struct {
	g    *errgroup.Group
	root string
}

func (wp walkerPipe) Pipe(ctx context.Context, _ <-chan any) <-chan any {
	out := make(chan any)

	wp.g.Go(func() error {
		defer close(out)

		return filepath.WalkDir(wp.root, func(
			name string,
			dirEntry fs.DirEntry,
			err error,
		) error {
			if err != nil {
				return err
			}
			if !dirEntry.Type().IsRegular() {
				return nil
			}

			select {
			case out <- name:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	})

	return out
}

type readPipe struct {
	g *errgroup.Group
}

type fileInfo struct {
	name string
	r    io.ReadCloser
}

func (rp readPipe) Pipe(ctx context.Context, in <-chan any) <-chan any {
	out := make(chan any)

	rp.g.Go(func() error {
		defer close(out)

		for v := range in {
			name := v.(string)
			f, err := os.Open(name)
			if err != nil {
				return err
			}

			select {
			case out <- fileInfo{name, f}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		return nil
	})

	return out
}

type digestPipe struct {
	g    *errgroup.Group
	algo Algorithm
}

type Checksum struct {
	Name string
	Hash []byte
}

func (dp digestPipe) Pipe(ctx context.Context, in <-chan any) <-chan any {
	out := make(chan any)

	dp.g.Go(func() error {
		defer close(out)

		for v := range in {
			file := v.(fileInfo)

			var hash hash.Hash
			switch dp.algo {
			case AlgorithmMD5:
				hash = md5.New()
			case AlgorithmSha1:
				hash = sha1.New()
			case AlgorithmSha256:
				hash = sha256.New()
			case AlgorithmSha512:
				hash = sha512.New()
			}

			if _, err := io.Copy(hash, file.r); err != nil {
				return err
			}

			select {
			case out <- Checksum{Name: file.name, Hash: hash.Sum(nil)}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		return nil
	})

	return out
}
