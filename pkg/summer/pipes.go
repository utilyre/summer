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
)

type Result[T any] struct {
	Val T
	Err error
}

type walkPipe struct {
	roots []string
}

func (wp walkPipe) Pipe(ctx context.Context, _ <-chan any) <-chan any {
	out := make(chan any)

	go func() {
		defer close(out)

		walk := func(name string, dirEntry fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("walker: %w", err)
			}
			if !dirEntry.Type().IsRegular() {
				return nil
			}

			select {
			case out <- Result[string]{Val: name}:
			case <-ctx.Done():
				return fmt.Errorf("walker: %w", ctx.Err())
			}
			return nil
		}

		for _, root := range wp.roots {
			if err := filepath.WalkDir(root, walk); err != nil {
				out <- Result[string]{Err: err}
			}
		}
	}()

	return out
}

type readPipe struct{}

type fileInfo struct {
	name string
	r    io.ReadCloser
}

func (rp readPipe) Pipe(ctx context.Context, in <-chan any) <-chan any {
	out := make(chan any)

	go func() {
		defer close(out)

		for v := range in {
			res := v.(Result[string])
			if res.Err != nil {
				out <- Result[fileInfo]{Err: res.Err}
				continue
			}

			f, err := os.Open(res.Val)
			if err != nil {
				out <- Result[fileInfo]{Err: fmt.Errorf("reader: %w", err)}
				continue
			}

			select {
			case out <- Result[fileInfo]{Val: fileInfo{res.Val, f}}:
			case <-ctx.Done():
				out <- Result[fileInfo]{Err: fmt.Errorf("reader: %w", ctx.Err())}
			}
		}
	}()

	return out
}

type digestPipe struct {
	algo Algorithm
}

// A Checksum represents name and hash of a particular file.
type Checksum struct {
	Name string
	Hash []byte
}

func (dp digestPipe) Pipe(ctx context.Context, in <-chan any) <-chan any {
	out := make(chan any)

	go func() {
		defer close(out)

		for v := range in {
			res := v.(Result[fileInfo])
			if res.Err != nil {
				out <- Result[Checksum]{Err: res.Err}
				continue
			}

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

			r := contextio.NewReader(ctx, res.Val.r)
			if _, err := io.Copy(hash, r); err != nil {
				out <- Result[Checksum]{Err: fmt.Errorf("digester: %w", err)}
				continue
			}

			select {
			case out <- Result[Checksum]{Val: Checksum{res.Val.name, hash.Sum(nil)}}:
			case <-ctx.Done():
				out <- Result[Checksum]{Err: fmt.Errorf("digester: %w", ctx.Err())}
			}
		}
	}()

	return out
}
