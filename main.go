package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "summer: not enough arguments")
		return
	}
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "summer: too many arguments")
		return
	}

	root := os.Args[1]
	sums, err := MD5All(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "summer: %s", err)
	}

	for path, sum := range sums {
		fmt.Printf("%x %s\n", sum, path)
	}
}

type Sum [md5.Size]byte

type result struct {
	path string
	sum  Sum
	err  error
}

func MD5All(root string) (map[string]Sum, error) {
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

	m := make(map[string]Sum)
	for r := range results {
		if r.err != nil {
			return nil, r.err
		}

		m[r.path] = r.sum
	}

	if err := <-errc; err != nil {
		return nil, err
	}

	return m, nil
}

const numDigesters int = 10

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
		case out <- &result{path: path, sum: md5.Sum(data), err: err}:
		case <-done:
			return
		}
	}
}
