package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/rozdolsky33/ocloud/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := cmd.Execute(ctx); err != nil {
		var serviceErr common.ServiceError
		if errors.As(err, &serviceErr) && serviceErr.GetHTTPStatusCode() == 401 {
			fmt.Fprintf(os.Stderr, "Authentication failed (%d %s). Please run `ocloud config session authenticate` to configure your profile \n", serviceErr.GetHTTPStatusCode(), serviceErr.GetCode())
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
