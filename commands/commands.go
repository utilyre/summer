package commands

import (
	"errors"
	"fmt"

	"gihtub.com/utilyre/summer/sum"
)

func Execute(args []string) error {
	if len(args) < 1 {
		return errors.New("not enough arguments")
	}
	if len(args) > 1 {
		return errors.New("too many arguments")
	}

	sums, err := sum.MD5All(args[0])
	if err != nil {
		return err
	}

	for path, sum := range sums {
		fmt.Printf("%x %s\n", sum, path)
	}

	return nil
}
