package image

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewImageCmd creates a new command for image-related operations
func NewImageCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "image",
		Aliases:       []string{"img"},
		Short:         "Explore OCI Compute images — list, get, and search",
		Long:          "List OCI Compute images in a compartment. Supports paging through large result sets and fuzzy search",
		Example:       "  ocloud compute image get\n  ocloud compute image list\n  ocloud compute image search <value>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))

	return cmd
}
