package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"gihtub.com/utilyre/summer/pkg/summer"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "summer: %v\n", err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("too few arguments")
	}
	if len(os.Args) > 2 {
		return errors.New("too many arguments")
	}

	checksums, err := summer.SumTree(context.Background(), os.Args[1])
	if err != nil {
		return err
	}

	for _, cs := range checksums {
		fmt.Printf("%x  %s\n", cs.Hash, cs.Name)
	}

	return nil
}
