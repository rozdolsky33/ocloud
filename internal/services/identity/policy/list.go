package policy

import (
	"context"
	"errors"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci/identity/policy"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// ListPolicies lists all policies in the specified compartment and prints their details in the specified format.
// ListPolicies lists available policies, presents them in an interactive TUI for selection, retrieves the selected policy, and prints it in either JSON or table form.
// It uses appCtx to initialize service clients and logging. If the user cancels the TUI flow the function returns nil; on failures to list, select, or fetch a policy it returns a wrapped error.
func ListPolicies(appCtx *app.ApplicationContext, useJSON bool, ocid string) error {
	ctx := context.Background()
	policyAdapter := policy.NewAdapter(appCtx.IdentityClient)
	service := NewService(policyAdapter, appCtx.Logger, ocid)
	policies, err := service.ListPolicies(ctx)

	if err != nil {
		return fmt.Errorf("listing policies: %w", err)
	}

	//TUI
	model := policy.NewPoliciesListModel(policies)
	id, err := tui.Run(model)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return fmt.Errorf("selecting policy: %w", err)
	}
	p, err := service.policyRepo.GetPolicy(ctx, id)
	if err != nil {
		return fmt.Errorf("getting policy: %w", err)
	}

	return PrintPolicyTable(p, appCtx, useJSON)
}
