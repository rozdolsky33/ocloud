package dynamicgroup

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci/identity/dynamicgroup"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

func ListDynamicGroups(appCtx *app.ApplicationContext, tenancyOCID string, useJSON bool) error {
	ctx := context.Background()
	dgAdapter := dynamicgroup.NewDynamicGroupAdapter(appCtx.IdentityClient, appCtx.Provider)
	service := NewService(dgAdapter, appCtx.Logger, tenancyOCID)

	dynamicGroups, err := service.dynamicGroupRepo.ListDynamicGroups(ctx, tenancyOCID)
	if err != nil {
		return fmt.Errorf("listing dynamic groups: %w", err)
	}

	// TUI — always show interactive picker, then output selected group
	model := dynamicgroup.NewDynamicGroupListModel(dynamicGroups)
	id, err := tui.Run(model)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting dynamic group: %w", err)
	}

	dg, err := service.dynamicGroupRepo.GetDynamicGroup(ctx, id)
	if err != nil {
		return fmt.Errorf("getting dynamic group: %w", err)
	}

	return PrintDynamicGroupInfo(dg, appCtx, useJSON)
}
