package main

import (
	"context"
	"os"
	"os/signal"
)

// sigintContext returns a context that is cancelled when CTRL+C is sent to the process.
func sigintContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal)
		signal.Notify(stop, os.Interrupt)
		<-stop

		cancel()
		signal.Stop(stop)
		close(stop)
	}()

	return ctx
}
