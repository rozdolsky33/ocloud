package instance

import (
	instaceFlags "github.com/rozdolsky33/ocloud/cmd/compute/flags"
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	"github.com/spf13/cobra"
)

// Long description for the list command
var listLong = `
List all instances in the specified compartment with pagination support.

This command displays information about running instances in the current compartment.
By default, it shows basic instance information such as name, ID, IP address, and shape.

The output is paginated, with a default limit of 20 instances per page. You can navigate
through pages using the --page flag and control the number of instances per page with
the --limit flag.

Additional Information:
- Use --image-details (-i) to include information about the image used by each instance
- Use --json (-j) to output the results in JSON format
- The command only shows running instances by default
`

// Examples for the list command
var listExamples = `
  # List all instances with default pagination (20 per page)
  ocloud compute instance list

  # List instances with custom pagination (10 per page, page 2)
  ocloud compute instance list --limit 10 --page 2

  # List instances and include image details
  ocloud compute instance list --image-details

  # List instances with image details (using shorthand flag)
  ocloud compute instance list -i

  # List instances and output in JSON format
  ocloud compute instance list --json

  # List instances with both image details and JSON output
  ocloud compute instance list --image-details --json
`

// NewListCmd creates a new command for listing instances
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all instances",
		Long:          listLong,
		Example:       listExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}

	// Add flags specific to the list command
	paginationFlags.LimitFlag.Add(cmd)
	paginationFlags.PageFlag.Add(cmd)
	instaceFlags.ImageDetailsFlag.Add(cmd)

	return cmd
}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	// Get pagination parameters
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	imageDetails := flags.GetBoolFlag(cmd, flags.FlagNameAllInformation, false)

	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(logger.CmdLogger, 1, "Running instance list command in", "compartment", appCtx.CompartmentName, "limit", limit, "page", page, "json", useJSON, "imageDetails", imageDetails)
	return instance.ListInstances(appCtx, limit, page, useJSON, imageDetails)
}
