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

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))

	return cmd
}
