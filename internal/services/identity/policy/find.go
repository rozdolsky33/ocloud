package policy

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci/identity/policy"
)

// FindPolicies searches for policies matching the given name pattern and prints their details in the specified format.
// It utilizes the application context for service initialization and handles output formatting via JSON or plain text.
func FindPolicies(appCtx *app.ApplicationContext, namePattern string, useJSON bool, ocid string) error {
	ctx := context.Background()
	policyAdapter := policy.NewAdapter(appCtx.IdentityClient)
	service := NewService(policyAdapter, appCtx.Logger, ocid)
	matchedPolicies, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding matched policies: %w", err)
	}

	return PrintPolicyInfo(matchedPolicies, appCtx, nil, useJSON)
}
