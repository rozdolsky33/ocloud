package util

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/spf13/cobra"
)

const (
	ScopeCompartment = "compartment"
	ScopeTenancy     = "tenancy"
)

// ResolveScope returns the final scope respecting precedence:
func ResolveScope(cmd *cobra.Command) string {
	if flags.GetBoolFlag(cmd, flags.FlagNameTenancyScope, false) {
		return ScopeTenancy
	}
	scope := strings.ToLower(flags.GetStringFlag(cmd, flags.FlagNameScope, ScopeCompartment))
	switch scope {
	case ScopeTenancy:
		return ScopeTenancy
	case ScopeCompartment, "":
		return ScopeCompartment
	default:
		return ScopeCompartment
	}
}

// ResolveParentID maps scope parent OCID
func ResolveParentID(scope string, appCtx *app.ApplicationContext) string {
	if scope == ScopeTenancy {
		return appCtx.TenancyID
	}
	if appCtx.CompartmentID != "" {
		return appCtx.CompartmentID
	}
	return appCtx.TenancyID
}
