package image

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/image"
	"github.com/spf13/cobra"
)

// Dedicated documentation for the list command (separate from get)
var listCmdLong = `
Interactively browse and search images in the specified compartment using a TUI.

This command launches a Bubble Tea-based terminal UI that loads available images and lets you:
- Search/filter images as you type
- Navigate the list
- Select a single image to view its details

After you pick an image, the tool prints detailed information about the selected image.
`

var listCmdExamples = `
  # Launch the interactive image browser
  ocloud compute image list

  # Use fuzzy search in the UI to quickly find what you need
  ocloud compute image list
`

// NewListCmd creates a new command for listing images
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List all images",
		Aliases:       []string{"l"},
		Long:          listCmdLong,
		Example:       listCmdExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}

	return cmd
}

// RunListCommand executes the interactive TUI image lister
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	ctx := cmd.Context()
	logger.LogWithLevel(logger.CmdLogger, 1, "Running image list (TUI) command in", "compartment", appCtx.CompartmentName)
	return image.ListImages(ctx, appCtx)
}
