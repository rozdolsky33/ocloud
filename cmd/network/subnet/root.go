package subnet

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewSubnetCmd creates a new command for subnet-related operations
func NewSubnetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "subnet",
		Aliases:       []string{"sub"},
		Short:         "Manage OCI Subnets",
		Long:          "Manage Oracle Cloud Infrastructure Subnets - list all subnets or find subnet by pattern.",
		Example:       "  ocloud network subnet list \n  ocloud network subnet find mysubnet",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))

	return cmd
}
