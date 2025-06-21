package instance

import (
	instaceFlags "github.com/rozdolsky33/ocloud/cmd/compute/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	"github.com/spf13/cobra"
)

// Long description for the find command
var findLong = `
Find instances in the specified compartment that match the given pattern.

The search is performed using a fuzzy matching algorithm that searches across multiple fields:

Searchable Fields:
- Name: Instance name
- ImageName: Name of the image used by the instance
- ImageOperatingSystem: Operating system of the image
- Tags: All instance tags in "key:value" format (e.g., "os_version:8.10")
- TagValues: Just the values of tags without keys (e.g., "8.10")

The search pattern is automatically wrapped with wildcards, so partial matches are supported.
For example, searching for "web" will match "webserver", "web-app", etc.

You can also search for specific tag values by using the tag key and value in your search pattern.
For example, "os_version:8.10" will find instances with that specific tag.
`

// Examples for the find command
var findExamples = `
  # Find instances with "web" in their name
  ocloud compute instance find web

  # Find instances with a specific tag value
  ocloud compute instance find os_version:8.10

  # Find instances with a specific tag value (searching just the value)
  ocloud compute instance find 8.10

  # Find instances with "api" in their name and include image details
  ocloud compute instance find api --image-details

  # Find instances with "server" in their name and output in JSON format
  ocloud compute instance find server --json

  # Find instances with "oracle" in their image operating system
  ocloud compute instance find oracle
`

// NewFindCmd creates a new command for finding instances by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [pattern]",
		Aliases:       []string{"f"},
		Short:         "Find instances by name pattern",
		Long:          findLong,
		Example:       findExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunFindCommand(cmd, args, appCtx)
		},
	}

	// Add flags specific to the find command
	instaceFlags.ImageDetailsFlag.Add(cmd)

	return cmd
}

// RunFindCommand handles the execution of the find command
func RunFindCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	imageDetails := flags.GetBoolFlag(cmd, flags.FlagNameImageDetails, false)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)

	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(logger.CmdLogger, 1, "Running instance find command", "pattern", namePattern, "in compartment", appCtx.CompartmentName, "json", useJSON)
	return instance.FindInstances(appCtx, namePattern, imageDetails, useJSON)
}
