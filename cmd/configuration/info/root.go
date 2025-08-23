package info

import (
	"github.com/spf13/cobra"
)

// Short description for the info command
var infoShort = "View information about ocloud environment configuration"

// Long description for the info command
var infoLong = `View information about ocloud environment configuration, such as tenancy mappings and other configuration details.

This command provides access to information about your ocloud environment, including tenancy mappings,
which allow you to associate tenancy names with their OCIDs and other metadata.`

// Examples for the info command
var infoExamples = `  ocloud config info map-file
  ocloud config i map-file --json
  ocloud config i map-file --realm OC1`

// NewInfoCmd creates a new cobra.Command for viewing information about ocloud environment configuration.
// It provides subcommands for viewing tenancy mapping information and other configuration details.
func NewInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "info",
		Aliases:       []string{"i"},
		Short:         infoShort,
		Long:          infoLong,
		Example:       infoExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(ViewMappingFile())

	return cmd
}
