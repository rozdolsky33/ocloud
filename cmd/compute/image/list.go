package image

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/image"
	"github.com/spf13/cobra"
)

// Dedicated documentation for the list command (separate from get)
var listLong = `
Interactively browse and search images in the specified compartment using a TUI.

This command launches terminal UI that loads available images and lets you:
- Search/filter image as you type
- Navigate the list
- Select a single image to view its details

After you pick an image, the tool prints detailed information about the selected image default table view or JSON format if specified with --json.
`

var listExamples = `
  # Launch the interactive images browser
  ocloud compute image list
  ocloud compute image list --json
`

// NewListCmd creates a new command for listing images
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List all images",
		Aliases:       []string{"l"},
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, appCtx)
		},
	}

	return cmd
}

// runListCommand executes the interactive TUI image lister
func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	ctx := cmd.Context()
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running image list (TUI) command in", "compartment", appCtx.CompartmentName)
	return image.ListImages(ctx, appCtx, useJSON)
}
