package summer

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/dolmen-go/contextio"
	"golang.org/x/sync/errgroup"
)

type walkPipe struct {
	g     *errgroup.Group
	roots []string
}

func (wp walkPipe) Pipe(ctx context.Context, _ <-chan any) <-chan any {
	out := make(chan any)

	wp.g.Go(func() error {
		defer close(out)

		walk := func(name string, dirEntry fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("walker: %w", err)
			}
			if !dirEntry.Type().IsRegular() {
				return nil
			}

			select {
			case out <- name:
			case <-ctx.Done():
				return fmt.Errorf("walker: %w", ctx.Err())
			}
			return nil
		}

		for _, root := range wp.roots {
			if err := filepath.WalkDir(root, walk); err != nil {
				return err
			}
		}

		return nil
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
				return fmt.Errorf("reader: %w", err)
			}

			select {
			case out <- fileInfo{name, f}:
			case <-ctx.Done():
				return fmt.Errorf("reader: %w", ctx.Err())
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

// A Checksum represents name and hash of a particular file.
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

			r := contextio.NewReader(ctx, file.r)
			if _, err := io.Copy(hash, r); err != nil {
				return fmt.Errorf("digester: %w", err)
			}

			select {
			case out <- Checksum{Name: file.name, Hash: hash.Sum(nil)}:
			case <-ctx.Done():
				return fmt.Errorf("digester: %w", ctx.Err())
			}
		}

		return nil
	})

	return out
}
