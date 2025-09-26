package instance

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ociInst "github.com/rozdolsky33/ocloud/internal/oci/compute/instance"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// ListInstances lists compute instances in the configured compartment, presents an interactive TUI to select one, and prints the selected instance's information.
// If useJSON is true the instance information is printed in JSON; otherwise it is formatted for human-readable output.
// It returns an error if creating the OCI clients, listing instances, selecting an instance via the TUI, or retrieving the selected instance fails; returning nil when the user cancels the selection.
func ListInstances(appCtx *app.ApplicationContext, useJSON bool) error {

	ctx := context.Background()

	computeClient, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	imageAdapter := ociInst.NewAdapter(computeClient, networkClient)
	service := NewService(imageAdapter, appCtx.Logger, appCtx.CompartmentID)
	allInstances, err := service.ListInstances(ctx)

	if err != nil {
		return fmt.Errorf("listing instances: %w", err)
	}

	//TUI
	model := ociInst.NewImageListModel(allInstances)
	id, err := tui.Run(model)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting instance: %w", err)
	}

	instance, err := service.instanceRepo.GetEnrichedInstance(ctx, id)

	if err != nil {
		return fmt.Errorf("getting instance: %w", err)
	}

	return PrintInstanceInfo(instance, appCtx, useJSON, true)
}
