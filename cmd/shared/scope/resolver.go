package scope

import (
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/spf13/cobra"
)

const (
	Compartment = "compartment"
	Tenancy     = "tenancy"
)

// ResolveScope returns the final scope respecting precedence:
func ResolveScope(cmd *cobra.Command) string {
	if flags.GetBoolFlag(cmd, flags.FlagNameTenancyScope, false) {
		return Tenancy
	}
	scope := strings.ToLower(flags.GetStringFlag(cmd, flags.FlagNameScope, Compartment))
	switch scope {
	case Tenancy:
		return Tenancy
	case Compartment, "":
		return Compartment
	default:
		return Compartment
	}
}

// ResolveParentID maps scope parent OCID
func ResolveParentID(scope string, appCtx *app.ApplicationContext) string {
	if scope == Tenancy {
		return appCtx.TenancyID
	}
	if appCtx.CompartmentID != "" {
		return appCtx.CompartmentID
	}
	return appCtx.TenancyID
}
