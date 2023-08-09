package sum

import (
	"crypto/md5"
	"fmt"
	"sync"

	"gihtub.com/utilyre/summer/fs"
)

const numDigesters int = 10

type Checksum [md5.Size]byte

type result struct {
	path string
	sum  Checksum
	err  error
}

func MD5All(root string) (map[string]Checksum, error) {
	done := make(chan struct{})
	defer close(done)

	paths, errc := fs.Walk(done, root)
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
