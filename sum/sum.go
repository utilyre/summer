package sum

import (
	"crypto/md5"

	"gihtub.com/utilyre/summer/channel"
)

type Sum [md5.Size]byte

func MD5All(root string) (map[string]Sum, error) {
	done := make(chan struct{})
	defer close(done)

	paths, errc := walk(done, root)
	contents := channel.Merge(
		read(done, paths),
		read(done, paths),
		read(done, paths),
		read(done, paths),
		read(done, paths),
	)
	sums := channel.Merge(
		digest(done, contents),
		digest(done, contents),
		digest(done, contents),
		digest(done, contents),
		digest(done, contents),
	)

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
