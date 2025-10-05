package image

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/image"
	"github.com/spf13/cobra"
)

// Long description for the find command
var findLong = `
FuzzySearch images in the specified compartment that match the given pattern.

The search is performed using a fuzzy matching algorithm that searches across multiple fields:

Searchable Fields:
- Name: Image name
- ImageOSVersion: Operating system version of the image
- OperatingSystem: Operating system of the image
- LunchMode: Launch mode of the image
- Tags: All image tags in format (e.g., "flock")

The search pattern is automatically wrapped with wildcards, so partial matches are supported.
For example, searching for "oracle" will match "oracle" etc.
`

// Examples for the find command
var findExamples = `
  # FuzzySearch images with "oracle" in their name
  ocloud compute image find oracle

  # FuzzySearch images with a specific operating system
  ocloud compute image find linux

  # FuzzySearch images with a specific tag value (searching just the value)
  ocloud compute image find 8.10

  # FuzzySearch images with "server" in their name and output in JSON format
  ocloud compute image find server --json

  # FuzzySearch images with a specific launch mode
  ocloud compute image find native
`

// NewFindCmd creates a new command for finding image by name pattern
func NewFindCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find [pattern]",
		Aliases:       []string{"f"},
		Short:         "FuzzySearch image by name pattern",
		Long:          findLong,
		Example:       findExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunFindCommand(cmd, args, appCtx)
		},
	}

	return cmd
}

// RunFindCommand handles the execution of the find command
func RunFindCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running image find command", "pattern", namePattern, "in compartment", appCtx.CompartmentName, "json", useJSON)
	return image.FindImages(appCtx, namePattern, useJSON)
}
