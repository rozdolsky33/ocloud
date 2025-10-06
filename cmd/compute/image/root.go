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
		Short:         "Manage OCI Compute images — list, paginate, and search.",
		Long:          "List OCI Compute images in a compartment. Supports paging through large result sets and filtering by value pattern.",
		Example:       "  ocloud compute image get\n  ocloud compute image list\n  ocloud compute image search <image-name>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewSearchCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))

	return cmd
}
