package image

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewImageCmd creates a new command for images-related operations
func NewImageCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "image",
		Aliases:       []string{"img"},
		Short:         "Manage OCI images",
		Long:          "Manage Oracle Cloud Infrastructure compute images - list all images or find image by name pattern.",
		Example:       "  ocloud compute image list\n  ocloud compute image find <images-name>",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))

	return cmd
}
