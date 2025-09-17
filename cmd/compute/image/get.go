package image

import (
	imageFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/image"
	"github.com/spf13/cobra"
)

// Dedicated documentation for the get command
var getLong = `
Get images in the specified compartment with pagination support.

This command retrieves available images in the current compartment.
By default, it shows basic image information such as name, ID, operating system, and launch mode.

The output is paginated, with a default limit of 20 images per page. You can navigate
through pages using the --page flag and control the number of images per page with
the --limit flag.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command shows all available images in the compartment
`

var getExamples = `
  # Get images with default pagination (20 per page)
  ocloud compute image get

  # Get images with custom pagination (10 per page, page 2)
  ocloud compute image get --limit 10 --page 2

  # Get images and output in JSON format
  ocloud compute image get --json

  # Get images with custom pagination and JSON output
  ocloud compute image get --limit 5 --page 3 --json
`

// NewGetCmd creates a new command for listing images
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Paginated Image Results",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGetCommand(cmd, appCtx)
		},
	}

	imageFlags.LimitFlag.Add(cmd)
	imageFlags.PageFlag.Add(cmd)

	return cmd
}

// RunGetCommand handles the execution of the list command
func RunGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, imageFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, imageFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running image list command in", "compartment", appCtx.CompartmentName, "limit", limit, "page", page, "json", useJSON)
	return image.GetImages(appCtx, limit, page, useJSON)
}
