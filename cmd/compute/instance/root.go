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
		Short:         "Explore OCI Compute instances — list, get, and search.",
		Long:          "List OCI Compute instances in a compartment. Supports paging through large result sets and fuzzy search",
		Example:       "  ocloud compute instance get\n  ocloud compute instance list\n  ocloud compute instance search <value>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))

	return cmd
}
