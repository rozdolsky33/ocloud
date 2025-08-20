package instance

import (
	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// NewInstanceCmd creates a new command for instance-related operations
func NewInstanceCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "instance",
		Aliases:       []string{"inst"},
		Short:         "Manage OCI Instances",
		Long:          "Manage Oracle Cloud Infrastructure Compute Instances - list all instances or find instances by name pattern.",
		Example:       "  ocloud compute instance list\n  ocloud compute instance find myinstance",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))

	return cmd
}
