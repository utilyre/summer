package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/utilyre/summer/pkg/summer"
)

var (
	algo          = summer.AlgorithmMD5
	timeout       int64
	readWorkers   int
	digestWorkers int
)

func init() {
	flag.Var(&algo, "algo", "sum using cryptographic hash function VALUE")
	flag.Int64Var(&timeout, "timeout", 0, "cancel operation after N milliseconds")
	flag.IntVar(&readWorkers, "read-workers", 1, "run N read workers in parallel")
	flag.IntVar(&digestWorkers, "digest-workers", 1, "run N digest workers in parallel")

	flag.Parse()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		defer cancel()
	}

	handleCancelSignals(cancel)

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "summer: %v\n", err)
	}
}

func run(ctx context.Context) error {
	roots := flag.Args()
	if len(roots) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		roots = append(roots, cwd)
	}

	checksums, err := summer.SumTree(
		ctx,
		roots,
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

func handleCancelSignals(cancel context.CancelFunc) {
	quitCh := make(chan os.Signal, 1)
	signal.Notify(
		quitCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGPIPE,
	)

	go func() {
		var lastSig os.Signal

		for sig := range quitCh {
			if lastSig != nil && lastSig.String() == sig.String() {
				os.Exit(1)
			}

			lastSig = sig
			cancel()
		}
	}()
}
