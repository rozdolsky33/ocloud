package gateway

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// GetGateway creates a network client from the provided application context and prepares a background context for gateway operations.
// It returns an error if creating the network client fails.
// The vcnName and useJSON parameters are currently unused.
func GetGateway(appCtx *app.ApplicationContext, vcnName string, useJSON bool) error {
	ctx := context.Background()
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}
	fmt.Println(networkClient)
	fmt.Println(ctx)
	return nil
}
