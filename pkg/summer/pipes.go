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

func walkDirs(ctx context.Context, roots []string) <-chan *Checksum {
	out := make(chan *Checksum)

	go func() {
		defer close(out)

		walk := func(name string, dirEntry fs.DirEntry, err error) error {
			cs := &Checksum{Name: name}

			if err != nil {
				cs.Err = fmt.Errorf("walk %s: %w", cs.Name, err)
				out <- cs
				return nil
			}
			if !dirEntry.Type().IsRegular() {
				return nil
			}

			select {
			case out <- cs:
			case <-ctx.Done():
				cs.Err = fmt.Errorf("walk %s: %w", cs.Name, ctx.Err())
				out <- cs

				// return non-nil error to abort the walk
				return cs.Err
			}

			return nil
		}

		for _, root := range roots {
			// NOTE: errors are handled by walk
			_ = filepath.WalkDir(root, walk)
		}
	}()

	return out
}

func walkFiles(ctx context.Context, names []string) <-chan *Checksum {
	out := make(chan *Checksum)

	go func() {
		defer close(out)

		for _, name := range names {
			cs := &Checksum{Name: name}

			info, err := os.Stat(name)
			if err != nil {
				cs.Err = fmt.Errorf("walk %s: %w", name, err)
				out <- cs
				continue
			}

			if info.IsDir() {
				cs.Err = fmt.Errorf("walk %s: is a directory", name)
				out <- cs
				continue
			}

			select {
			case out <- cs:
			case <-ctx.Done():
				cs.Err = fmt.Errorf("walk %s: %w", cs.Name, ctx.Err())
				out <- cs
			}
		}
	}()

	return out
}

type readPipe struct{}

func (rp readPipe) Pipe(ctx context.Context, in <-chan *Checksum) <-chan *Checksum {
	out := make(chan *Checksum)

	go func() {
		defer close(out)

		for cs := range in {
			if cs.Err != nil {
				out <- cs
				continue
			}

			f, err := os.Open(cs.Name)
			if err != nil {
				cs.Err = fmt.Errorf("read %s: %w", cs.Name, err)
				out <- cs
				continue
			}

			cs.body = f

			select {
			case out <- cs:
			case <-ctx.Done():
				cs.Err = fmt.Errorf("read %s: %w", cs.Name, ctx.Err())
				out <- cs
			}
		}
	}()

	return out
}

type digestPipe struct {
	algo Algorithm
}

func (dp digestPipe) Pipe(ctx context.Context, in <-chan *Checksum) <-chan *Checksum {
	out := make(chan *Checksum)

	go func() {
		defer close(out)

		for cs := range in {
			if cs.Err != nil {
				out <- cs
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

			r := contextio.NewReader(ctx, cs.body)
			if _, err := io.Copy(hash, r); err != nil {
				cs.Err = fmt.Errorf("digest %s: %w", cs.Name, err)
				continue
			}

			cs.Hash = hash.Sum(nil)

			select {
			case out <- cs:
			case <-ctx.Done():
				cs.Err = fmt.Errorf("digest %s: %w", cs.Name, ctx.Err())
				out <- cs
			}
		}
	}()

	return out
}
