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

	results, err := summer.SumTree(
		ctx,
		args,
		summer.WithAlgorithm(algo),
		summer.WithReadJobs(readJobs),
		summer.WithDigestJobs(digestJobs),
	)
	if err != nil {
		return err
	}

	for _, res := range results {
		if res.Err != nil {
			fmt.Fprintln(os.Stderr, "summer:", res.Err)
			continue
		}

		fmt.Printf("%x  %s\n", res.Val.Hash, res.Val.Name)
	}

	return nil
}
