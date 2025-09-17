package compartment

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci/identity/compartment"
	"github.com/rozdolsky33/ocloud/internal/tui/listx"
)

func ListCompartments(appCtx *app.ApplicationContext, ocid string, useJSON bool) error {
	ctx := context.Background()
	compartmentAdapter := compartment.NewCompartmentAdapter(appCtx.IdentityClient, ocid)
	service := NewService(compartmentAdapter, appCtx.Logger, ocid)

	compartments, err := service.compartmentRepo.ListCompartments(ctx, ocid)
	if err != nil {
		return fmt.Errorf("getting compartment: %w", err)
	}

	//TUI
	model := compartment.NewPoliciesListModel(compartments)
	id, err := listx.Run(model)
	if err != nil {
		if errors.Is(err, listx.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting compartment: %w", err)
	}
	c, err := service.compartmentRepo.GetCompartment(ctx, id)
	if err != nil {
		return fmt.Errorf("getting compartment: %w", err)
	}

	return PrintCompartmentInfo(c, appCtx, useJSON)
}
