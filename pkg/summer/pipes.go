package summer

import (
	"context"
	"crypto/md5"
	"io/fs"
	"os"
	"path/filepath"
)

type walkerPipe struct {
	root string
}

func (wp walkerPipe) Pipe(ctx context.Context, _ <-chan any) <-chan any {
	out := make(chan any)

	go func() {
		defer close(out)

		// TODO: handle error
		filepath.WalkDir(wp.root, func(
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
	}()

	return out
}

type readerPipe struct{}

type fileInfo struct {
	name string
	data []byte
}

func (rp readerPipe) Pipe(ctx context.Context, in <-chan any) <-chan any {
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
			case out <- fileInfo{name, data}:
			case <-ctx.Done():
				panic("TODO")
			}
		}
	}()

	return out
}

type digesterPipe struct{}

type ChecksumInfo struct {
	Name     string
	Checksum []byte
}

func (dp digesterPipe) Pipe(ctx context.Context, in <-chan any) <-chan any {
	out := make(chan any)

	go func() {
		defer close(out)

		for v := range in {
			file := v.(fileInfo)
			sum := md5.Sum(file.data)

			select {
			case out <- ChecksumInfo{Name: file.name, Checksum: sum[:]}:
			case <-ctx.Done():
				panic("TODO")
			}
		}
	}()

	return out
}
