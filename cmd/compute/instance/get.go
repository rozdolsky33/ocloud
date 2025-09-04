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

var getLong = `
Get all instances in the specified compartment with pagination support.

This command displays information about running instances in the current compartment.
By default, it shows basic instance information such as name, ID, IP address, and shape.

The output is paginated, with a default limit of 20 instances per page. You can navigate
through pages using the --page flag and control the number of instances per page with
the --limit flag.

Additional Information:
- Use --all (-A) to include detailed information about the instance
- Use --json (-j) to output the results in JSON format
- The command only shows running instances by default
`

var getExamples = `
  # Get all instances with default pagination (20 per page)
  ocloud compute instance get

  # Get instances with custom pagination (10 per page, page 2)
  ocloud compute instance get --limit 10 --page 2

  # Get instances and include instance details
  ocloud compute instance get --all

  # Get instances with instance details (using shorthand flag)
  ocloud compute instance get -A

  # Get instances and output in JSON format
  ocloud compute instance get --json

  # Get instances with both instance details and JSON output
  ocloud compute instance get --all --json
`

// NewGetCmd creates a new command for listing instances
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Paginated Instance Results",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGetCommand(cmd, appCtx)
		},
	}

	paginationFlags.LimitFlag.Add(cmd)
	paginationFlags.PageFlag.Add(cmd)
	instaceFlags.AllInfoFlag.Add(cmd)

	return cmd
}

// RunGetCommand handles the execution of the list command
func RunGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	imageDetails := flags.GetBoolFlag(cmd, flags.FlagNameAllInformation, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running instance get command in", "compartment", appCtx.CompartmentName, "limit", limit, "page", page, "json", useJSON, "imageDetails", imageDetails)
	return instance.GetInstances(appCtx, useJSON, limit, page, imageDetails)
}
