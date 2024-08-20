package main

import (
	"github.com/spf13/cobra"
	"github.com/utilyre/summer/pkg/summer"
)

var (
	timeout int64
	algo    = summer.AlgorithmMD5
)

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "summer",
		Short:         "High-performance utility for generating checksums in parallel",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	fset := cmd.PersistentFlags()
	fset.Int64Var(&timeout, "timeout", 0, "cancel operation after given milliseconds")
	fset.Var(&algo, "algo", "sum using given cryptographic hash function")

	cmdGenerate := &cobra.Command{
		Use:   "generate [files]",
		Short: "Generate checksums for given files",
		RunE:  generate,
	}
	fset = cmdGenerate.Flags()
	fset.Int("read-jobs", 1, "run given number of read jobs in parallel")
	fset.Int("digest-jobs", 1, "run given number of digest jobs in parallel")

	cmd.AddCommand(cmdGenerate)
	return cmd
}
