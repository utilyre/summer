package command

import (
	"errors"
	"fmt"

	"gihtub.com/utilyre/summer/sum"
)

func Execute(args []string) error {
	if len(args) > 1 {
		return errors.New("too many arguments")
	}

	root := "."
	if len(args) == 1 {
		root = args[0]
	}

	sums, err := sum.MD5All(root)
	if err != nil {
		return err
	}

	for path, sum := range sums {
		fmt.Printf("%x %s\n", sum, path)
	}

	return nil
}
