package version

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/buildinfo"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/spf13/cobra"
)

// VersionInfo encapsulates the version command functionality
// It wraps a cobra.Command and provides methods to handle version information display
type VersionInfo struct {
	cmd *cobra.Command
}

// NewVersionCommand creates and configures a new version command
// Returns a *cobra.Command that can be added to the root command
// This function was refactored to return *cobra.Command directly instead of *VersionInfo
// to fix an issue with adding the command to the root command
func NewVersionCommand() *cobra.Command {
	vc := &VersionInfo{}

	vc.cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  "Print the version, build time, and git commit hash information",
		RunE:  vc.runCommand,
	}

	return vc.cmd
}

// runCommand handles the main command execution
func (vc *VersionInfo) runCommand(cmd *cobra.Command, args []string) error {
	return vc.printVersionInfo()
}

// printVersionInfo displays the version information
func (vc *VersionInfo) printVersionInfo() error {
	PrintVersionInfo()
	return nil
}

// PrintVersionInfo prints complete version information to stdout
// This function was updated to print all version information (version, commit hash, and build time)
// to ensure consistency between the version command and the version flag
func PrintVersionInfo() {
	fmt.Printf("Version:    %s\n", buildinfo.Version)
	fmt.Printf("Commit:     %s\n", buildinfo.CommitHash)
	fmt.Printf("Built:      %s\n", buildinfo.BuildTime)
}

// PrintVersion prints version information to stdout
// This function is used by the root command when the --version flag is specified
// It was added to fix an issue where cmd/root.go was calling version.PrintVersion()
// which didn't exist in the version package
func PrintVersion() {
	PrintVersionInfo()
}

// AddVersionFlag adds a version flag to the root command
// This function adds a global persistent flag to support the --version/-v flag
// and sets up a PersistentPreRunE hook to check for the flag and print version information
// Note: This function preserves any existing PersistentPreRunE hook by storing it and
// calling it after checking for the version flag
func AddVersionFlag(rootCmd *cobra.Command) {
	// Register a global persistent flag to support short form (e.g., `ocloud -v`)
	rootCmd.PersistentFlags().BoolP(flags.FlagNameVersion, flags.FlagShortVersion, false, flags.FlagDescVersion)

	// Store the original PersistentPreRunE function
	originalPreRun := rootCmd.PersistentPreRunE

	// Override the persistent pre-run hook to check for the `-v` flag
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if versionFlag := flags.GetBoolFlag(cmd, flags.FlagNameVersion, false); versionFlag {
			PrintVersionInfo()
			return nil
		}

		// Call the original PersistentPreRunE if it exists
		if originalPreRun != nil {
			return originalPreRun(cmd, args)
		}
		return nil
	}
}
