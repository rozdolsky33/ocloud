package dynamicgroup

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewDynamicGroupCmd creates a new command for dynamic-group-related operations
func NewDynamicGroupCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "dynamic-group",
		Aliases:       []string{"dynamicgroup", "dg"},
		Short:         "Explore OCI Dynamic Groups",
		Long:          "Explore Oracle Cloud Infrastructure Dynamic Groups: list, get and search",
		Example:       "  ocloud identity dynamic-group list \n  ocloud identity dynamic-group get <ocid> \n  ocloud identity dynamic-group search <value>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))

	return cmd
}
