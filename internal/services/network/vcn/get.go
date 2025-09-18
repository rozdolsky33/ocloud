package vcn

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociVcn "github.com/rozdolsky33/ocloud/internal/oci/network/vcn"
)

// GetVCN retrieves a VCN by OCID and prints its summary or JSON.
func GetVCN(appCtx *app.ApplicationContext, vcnID string, useJSON bool) error {
	if vcnID == "" {
		return fmt.Errorf("vcn OCID is required")
	}

	client, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}
	adapter := ociVcn.NewAdapter(client)
	service := NewService(adapter, appCtx.Logger, appCtx.CompartmentID)

	v, err := service.GetVcn(context.Background(), vcnID)
	if err != nil {
		return fmt.Errorf("getting vcn: %w", err)
	}

	return PrintVCNSummary(ToDTO(v), appCtx, useJSON)
}
