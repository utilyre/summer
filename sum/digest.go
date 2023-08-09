package sum

import (
	"crypto/md5"
	"fmt"
	"os"
)

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

const numDigesters int = 10

func digseter(done <-chan struct{}, paths <-chan string, out chan<- *result) {
	for path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			err = &DigestError{path: path, err: err}
		}

		select {
		case out <- &result{
			path: path,
			sum:  md5.Sum(data),
			err:  err,
		}:
		case <-done:
			return
		}
	}
}
