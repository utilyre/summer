package commands

import (
	"errors"
	"fmt"
	"os"

	"gihtub.com/utilyre/summer/sum"
)

func Execute(args []string) error {
	if len(args) > 1 {
		return errors.New("too many arguments")
	}

	root := ""
	if len(args) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		root = cwd
	} else {
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
