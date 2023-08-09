package sum

import (
	"crypto/md5"
)

type Sum [md5.Size]byte

func MD5All(root string) (map[string]Sum, error) {
	done := make(chan struct{})
	defer close(done)

	paths, errc := walk(done, root)
	contents := read(done, paths)
	sums := digest(done, contents)

	m := make(map[string]Sum)
	for sum := range sums {
		if sum.err != nil {
			continue
		}

		m[sum.data.path] = sum.data.sum
	}

	if err := <-errc; err != nil {
		return nil, err
	}

	return m, nil
}
