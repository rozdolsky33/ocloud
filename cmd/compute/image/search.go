package image

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	cfgflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/image"
	"github.com/spf13/cobra"
)

var searchLong = `
Search for images in the specified compartment that match the given pattern.

The search uses a fuzzy, prefix, and substring matching algorithm across many indexed fields.
You can search using any of the following fields (partial matches are supported):

Searchable fields:
- OCID: Image OCID
- Name: Display name of the image
- OSVersion: Operating system version
- LaunchMode: Launch mode of the image
- OperatingSystem: Operating system of the image

The search pattern is case-insensitive. For very specific inputs (like full OCID),
the search first tries exact and substring matches; otherwise it falls back to broader fuzzy search.
`

var searchExamples = `
  # Search by display name (substring)
  ocloud compute image search ubuntu

  # Search by OS
  ocloud compute image search "Oracle-Linux"

  # Search by OS version
  ocloud compute image search 8.10

  # Show more details in the output
  ocloud compute image search api --all

  # Output in JSON format
  ocloud compute image search server --json
`

// NewSearchCmd creates a new command for finding images by name pattern
func NewSearchCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search [pattern]",
		Aliases:       []string{"s"},
		Short:         "Search images by name pattern",
		Long:          searchLong,
		Example:       searchExamples,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearchCommand(cmd, args, appCtx)
		},
	}

	return cmd
}

// runSearchCommand handles the execution of the search command
func runSearchCommand(cmd *cobra.Command, args []string, appCtx *app.ApplicationContext) error {
	namePattern := args[0]
	useJSON := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameJSON, false)
	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running image search command", "pattern", namePattern, "in compartment", appCtx.CompartmentName, "json", useJSON)
	return image.SearchImages(appCtx, namePattern, useJSON)
}
