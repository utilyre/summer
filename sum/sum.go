package sum

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type Checksum [md5.Size]byte

func MD5All(root string) (map[string]Checksum, error) {
	done := make(chan struct{})
	defer close(done)

	paths, errc := walk(done, root)
	results := make(chan *result)

	var wg sync.WaitGroup
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			digseter(done, paths, results)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	m := make(map[string]Checksum)
	for r := range results {
		if r.err != nil {
			return nil, fmt.Errorf("sum: %w", r.err)
		}

		m[r.path] = r.sum
	}

	if err := <-errc; err != nil {
		return nil, fmt.Errorf("sum: %w", err)
	}

	return m, nil
}

const numDigesters int = 10

type DigestError struct {
	path string
	err  error
}

func (e *DigestError) Error() string {
	return fmt.Sprintf("digesting %s failed due to %s", e.path, e.err)
}

func (e *DigestError) Unwrap() error {
	return e.err
}

type result struct {
	path string
	sum  Checksum
	err  error
}

func walk(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(paths)

		errc <- filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() {
				return nil
			}

			select {
			case paths <- path:
			case <-done:
				return errors.New("walk cancelled")
			}
			return nil
		})
	}()

	return paths, errc
}

func digseter(done <-chan struct{}, paths <-chan string, out chan<- *result) {
	for path := range paths {
		data, err := os.ReadFile(path)

		select {
		case out <- &result{
			path: path,
			sum:  md5.Sum(data),
			err:  &DigestError{path: path, err: err},
		}:
		case <-done:
			return
		}
	}
}
