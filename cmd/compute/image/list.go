package image

import (
	imageFlags "github.com/rozdolsky33/ocloud/cmd/compute/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/images"
	"github.com/spf13/cobra"
)

// NewListCmd creates a new command for listing images
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all images",
		Long:          "List all images in the specified compartment with pagination support.",
		Example:       "  ocloud compute image list\n  ocloud compute image list --limit 10 --page 2\n  ocloud compute image list --json\n  ocloud compute image list \n  ocloud compute image list",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}

	// Add flags specific to the list command
	imageFlags.LimitFlag.Add(cmd)
	imageFlags.PageFlag.Add(cmd)

	return cmd
}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {

	// Get pagination parameters
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, imageFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, imageFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)

	logger.LogWithLevel(logger.CmdLogger, 1, "Running image list command in", "compartment", appCtx.CompartmentName, "limit", limit, "page", page, "json", useJSON, "imageDetails")
	return images.ListImages(appCtx, limit, page, useJSON)
}
