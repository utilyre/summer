package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/utilyre/summer/internal/cli"
	"github.com/utilyre/summer/pkg/summer"
)

func runGenerate(cmd *cobra.Command, args []string) error {
	ctx, cancel := cli.NewContext(time.Duration(timeout) * time.Millisecond)
	defer cancel()

	if len(args) == 0 {
		args = append(args, ".")
	}

	fset := cmd.LocalFlags()
	readJobs, err := fset.GetInt("read-jobs")
	if err != nil {
		return err
	}
	digestJobs, err := fset.GetInt("digest-jobs")
	if err != nil {
		return err
	}

	checksums, err := summer.SumTree(
		ctx,
		args,
		summer.WithAlgorithm(algo),
		summer.WithReadJobs(readJobs),
		summer.WithDigestJobs(digestJobs),
	)
	if err != nil {
		return err
	}

	for _, cs := range checksums {
		fmt.Printf("%x  %s\n", cs.Hash, cs.Name)
	}

	return nil
}
