package sum

import (
	"context"
	"crypto/md5"

	"gihtub.com/utilyre/summer/channel"
)

type Sum [md5.Size]byte

func MD5All(root string) (map[string]Sum, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	paths, errc := walk(ctx, root)
	contents := channel.Merge(
		read(ctx, paths),
		read(ctx, paths),
		read(ctx, paths),
		read(ctx, paths),
		read(ctx, paths),
	)
	sums := channel.Merge(
		digest(ctx, contents),
		digest(ctx, contents),
		digest(ctx, contents),
		digest(ctx, contents),
		digest(ctx, contents),
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
