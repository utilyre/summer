package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/utilyre/summer/internal/cli"
	"github.com/utilyre/summer/pkg/summer"
)

func runGenerate(cmd *cobra.Command, args []string) error {
	ctx, cancel := cli.NewContext(time.Duration(timeout) * time.Millisecond)
	defer cancel()

	fset := cmd.LocalFlags()
	openFileJobs, err := fset.GetInt("open-file-jobs")
	if err != nil {
		return err
	}
	digestJobs, err := fset.GetInt("digest-jobs")
	if err != nil {
		return err
	}
	recursive, err := fset.GetBool("recursive")
	if err != nil {
		return err
	}

	s, err := summer.New(
		summer.WithAlgorithm(algo),
		summer.WithOpenFileJobs(openFileJobs),
		summer.WithDigestJobs(digestJobs),
		summer.WithRecursive(recursive),
	)
	if err != nil {
		return err
	}

	checksums, err := s.Sum(ctx, args...)
	if err != nil {
		return err
	}

	for cs := range checksums {
		if cs.Err != nil {
			fmt.Fprintln(os.Stderr, "summer:", cs.Err)
			continue
		}

		fmt.Printf("%x  %s\n", cs.Hash, cs.Name)
	}

	return nil
}
