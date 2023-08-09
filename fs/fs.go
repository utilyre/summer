package fs

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
)

type Error struct {
	path string
	err  error
}

func (e *Error) Error() string {
	return fmt.Sprintf("yielding %s failed since %s", e.path, e.err)
}

func (e *Error) Unwrap() error {
	return e.err
}

func Walk(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(paths)

		errc <- filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return &Error{path: path, err: err}
			}
			if !d.Type().IsRegular() {
				return nil
			}

			select {
			case paths <- path:
			case <-done:
				return &Error{path: path, err: errors.New("walk cancelled")}
			}
			return nil
		})
	}()

	return paths, errc
}
