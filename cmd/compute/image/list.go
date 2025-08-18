package image

import (
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/image"
	"github.com/spf13/cobra"
)

var listLong = `
List all images in the specified compartment with pagination support.

This command displays information about available images in the current compartment.
By default, it shows basic image information such as name, ID, operating system, and launch mode.

The output is paginated, with a default limit of 20 images per page. You can navigate
through pages using the --page flag and control the number of images per page with
the --limit flag.

Additional Information:
- Use --json (-j) to output the results in JSON format
- The command shows all available images in the compartment
`

var listExamples = `
  # List all images with default pagination (20 per page)
  ocloud compute image list

  # List images with custom pagination (10 per page, page 2)
  ocloud compute image list --limit 10 --page 2

  # List images and output in JSON format
  ocloud compute image list --json

  # List images with custom pagination and JSON output
  ocloud compute image list --limit 5 --page 3 --json
`

// NewListCmd creates a new command for listing images
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all images",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}

	paginationFlags.LimitFlag.Add(cmd)
	paginationFlags.PageFlag.Add(cmd)

	return cmd
}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {

	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)

	logger.LogWithLevel(logger.CmdLogger, 1, "Running image list command in", "compartment", appCtx.CompartmentName, "limit", limit, "page", page, "json", useJSON)
	return image.ListImages(appCtx, limit, page, useJSON)
}
