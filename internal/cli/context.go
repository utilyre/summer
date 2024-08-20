package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	var cancel1, cancel2 context.CancelFunc
	ctx, cancel1 := context.WithCancel(context.Background())
	if timeout > 0 {
		ctx, cancel2 = context.WithTimeout(ctx, timeout)
	}

	cancel := func() {
		if cancel2 != nil {
			cancel2()
		}

		cancel1()
	}

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
		wasSIGINT := false

		for sig := range quitCh {
			if wasSIGINT && sig == syscall.SIGINT {
				os.Exit(1)
			}

			wasSIGINT = sig == syscall.SIGINT
			cancel()
		}
	}()

	return ctx, cancel
}
