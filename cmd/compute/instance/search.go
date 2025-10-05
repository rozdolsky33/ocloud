package instance

import (
	instaceFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	"github.com/spf13/cobra"
)

var searchLong = `
Search for instances in the specified compartment that match the given pattern.

The search is performed using a fuzzy matching algorithm that searches across multiple fields:

Searchable Fields:
- Name: Instance name
- InstanceName: Name of the instance used by the instance
- InstanceOperatingSystem: Operating system of the instance
- TagValues: Just the values of tags without keys (e.g., "8.10")

The search pattern is automatically wrapped with wildcards, so partial matches are supported.
For example, searching for "web" will match "webserver" etc.
`

var searchExamples = `
  # Search for instances with "web" in their name
  ocloud compute instance search web

  # Search for instances with a specific tag value (searching just the value)
  ocloud compute instance search 8.10

  # Search for instances with "api" in their name and include instance details
  ocloud compute instance search api --all

  # Search for instances with "server" in their name and output in JSON format
  ocloud compute instance search server --json

  # Search for instances with "oracle" in their instance operating system
  ocloud compute instance search oracle
`

// NewSearchCmd creates a new command for finding instances by name pattern
func NewSearchCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search [pattern]",
		Aliases:       []string{"s"},
		Short:         "Search instances by name pattern",
		Long:          searchLong,
		Example:       searchExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearchCommand(cmd, args, appCtx)
		},
	}

	instaceFlags.AllInfoFlag.Add(cmd)

	return cmd
}

// RunSearchCommand handles the execution of the find command
func runSearchCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	showDetails := flags.GetBoolFlag(cmd, flags.FlagNameAll, false)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running instance search command", "pattern", namePattern, "in compartment", appCtx.CompartmentName, "json", useJSON)
	return instance.SearchInstances(appCtx, namePattern, useJSON, showDetails)
}
