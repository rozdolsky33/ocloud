package compartment

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci/identity/compartment"
)

// FindCompartments searches and displays compartments matching a given name pattern.
func FindCompartments(appCtx *app.ApplicationContext, namePattern string, useJSON bool, ocid string) error {
	ctx := context.Background()
	// Create the infrastructure adapter.
	compartmentAdapter := compartment.NewCompartmentAdapter(appCtx.IdentityClient, ocid)

	// Create the application service, injecting the adapter.
	service := NewService(compartmentAdapter, appCtx.Logger, ocid)

	matchedCompartments, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding matched compartments: %w", err)
	}

	return PrintCompartmentsInfo(matchedCompartments, appCtx, nil, useJSON)
}
