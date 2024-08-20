package main

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed version.txt
var version string

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Println(version)
	return nil
}
