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
		Short:         "Manage OCI Image",
		Long:          "Manage Oracle Cloud Infrastructure Compute Image - list all image or find image by name pattern.",
		Example:       "  ocloud compute image list\n  ocloud compute image find <image-name>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewGetCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))
	cmd.AddCommand(NewListCmd(appCtx))

	return cmd
}
