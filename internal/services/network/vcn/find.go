package vcn

//import (
//	"context"
//	"fmt"
//
//	"github.com/rozdolsky33/ocloud/internal/app"
//	"github.com/rozdolsky33/ocloud/internal/oci"
//	oci_vcn "github.com/rozdolsky33/ocloud/internal/oci/network/vcn"
//)
//
//// FindVCNs searches for VCNs by pattern in the current compartment and prints a table or JSON.
//func FindVCNs(appCtx *app.ApplicationContext, pattern string, useJSON bool) error {
//	if pattern == "" {
//		return fmt.Errorf("pattern is required")
//	}
//
//	client, err := oci.NewNetworkClient(appCtx.Provider)
//	if err != nil {
//		return fmt.Errorf("creating network client: %w", err)
//	}
//	adapter := oci_vcn.NewAdapter(client)
//	service := NewService(adapter, appCtx.Logger, appCtx.CompartmentID)
//
//	vcns, err := service.Find(context.Background(), pattern)
//	if err != nil {
//		return fmt.Errorf("finding vcns: %w", err)
//	}
//
//	var dtos []*VCNDTO
//	for _, v := range vcns {
//		dtos = append(dtos, ToDTO(v))
//	}
//
//	return PrintVCNsTable(dtos, appCtx, useJSON)
//}
