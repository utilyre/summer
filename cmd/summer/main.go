package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gihtub.com/utilyre/summer/pkg/summer"
)

var (
	algo          = summer.AlgorithmMD5
	readWorkers   int
	digestWorkers int
)

func init() {
	flag.Var(&algo, "algo", "sum using cryptographic hash function VALUE")
	flag.IntVar(&readWorkers, "read-workers", 1, "run N read workers in parallel")
	flag.IntVar(&digestWorkers, "digest-workers", 1, "run N digest workers in parallel")

	flag.Parse()
}

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGPIPE,
	)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "summer: %v\n", err)
	}
}

func run(ctx context.Context) error {
	if flag.NArg() != 1 {
		return errors.New("too many or too few arguments")
	}

	checksums, err := summer.SumTree(
		ctx,
		flag.Arg(0),
		summer.WithAlgorithm(algo),
		summer.WithReadWorkers(readWorkers),
		summer.WithDigestWorkers(digestWorkers),
	)
	if err != nil {
		return err
	}

	for _, cs := range checksums {
		fmt.Printf("%x  %s\n", cs.Hash, cs.Name)
	}

	return nil
}
