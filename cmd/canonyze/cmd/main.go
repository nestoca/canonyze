package main

import (
	"context"
	"os"
	"os/signal"
)

var version string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		stop()
	}()

	cmd := NewCmd()
	if err := cmd.ExecuteContext(ctx); err != nil {
		os.Exit(-1)
	}
}
