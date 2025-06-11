package main

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/cmd"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	interrupted := make(chan os.Signal, 1)
	defer close(interrupted)
	signal.Notify(interrupted, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupted)

	ctx, cancel := context.WithCancelCause(context.Background())

	go func() {
		select {
		case <-interrupted:
			cancel(fmt.Errorf("command interrupted"))
		}
	}()

	cmd.Execute(ctx)
}
