package gateway

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

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
