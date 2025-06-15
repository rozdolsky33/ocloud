package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

// This will be populated dynamically at build time
var version = "unknown"

// NewVersionCmd creates a new version command
func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of ocloud CLI",
		Long:  "Print the version number of ocloud CLI",
		Run: func(cmd *cobra.Command, args []string) {
			PrintVersion()
		},
	}

	return versionCmd
}

// PrintVersion prints the version information
func PrintVersion() {
	fmt.Printf("ocloud version: %s\n", version)
}

// AddVersionFlag adds a version flag to the root command
func AddVersionFlag(rootCmd *cobra.Command) {
	// Register a global persistent flag to support short form (e.g., `ocloud -v`)
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version number of ocloud CLI")

	// Store the original PersistentPreRunE function
	originalPreRun := rootCmd.PersistentPreRunE

	// Override the persistent pre-run hook to check for the `-v` flag
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			PrintVersion()
			return nil
		}

		// Call the original PersistentPreRunE if it exists
		if originalPreRun != nil {
			return originalPreRun(cmd, args)
		}
		return nil
	}
}
