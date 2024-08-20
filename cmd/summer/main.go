package main

import (
	"fmt"
	"os"
)

func main() {
	cmd := newCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stdout, "%s: %v\n", cmd.Name(), err)
		os.Exit(1)
	}
}
