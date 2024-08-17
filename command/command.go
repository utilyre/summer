package command

import (
	"context"
	"errors"
	"fmt"

	"gihtub.com/utilyre/summer/pkg/summer"
)

func Execute(args []string) error {
	if len(args) > 1 {
		return errors.New("too many arguments")
	}

	root := "."
	if len(args) == 1 {
		root = args[0]
	}

	sums, err := summer.SumTree(context.Background(), root, summer.AlgorithmMD5)
	if err != nil {
		return err
	}

	for _, info := range sums {
		fmt.Printf("%x  %s\n", info.Checksum, info.Name)
	}

	return nil
}
