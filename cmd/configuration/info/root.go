package info

import (
	"github.com/spf13/cobra"
)

// NewInfoCmd creates a new cobra.Command for viewing information about ocloud environment configuration.
// It provides subcommands for viewing tenancy mapping information and other configuration details.
func NewInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "info",
		Aliases:       []string{"i"},
		Short:         "View information about ocloud environment configuration",
		Long:          "View information about ocloud environment configuration, such as tenancy mappings and other configuration details.",
		Example:       "  ocloud config info map-file\n  ocloud config info map-file --json\n  ocloud config info map-file --realm OC1",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(ViewMappingFile())

	return cmd
}
