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

	errcs := []<-chan error{}

	pathc, errc := walk(ctx, root)
	errcs = append(errcs, errc)

	filecs := []<-chan file{}
	for i := 0; i < 5; i++ {
		c, errc := read(ctx, pathc)
		filecs = append(filecs, c)
		errcs = append(errcs, errc)
	}
	filec := channel.Merge(filecs...)

	checksumcs := []<-chan checksum{}
	for i := 0; i < 5; i++ {
		c, errc := digest(ctx, filec)
		checksumcs = append(checksumcs, c)
		errcs = append(errcs, errc)
	}
	checksumc := channel.Merge(checksumcs...)
	errc = channel.Merge(errcs...)

	m := make(map[string]Sum)
	for checksum := range checksumc {
		m[checksum.path] = checksum.sum
	}

	for err := range errc {
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
