package instance

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	"github.com/spf13/cobra"
)

// NewListCmd creates a new command for listing instances
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all instances",
		Long:          "List all instances in the specified compartment with pagination support.",
		Example:       "  ocloud compute instance list\n  ocloud compute instance list --limit 10 --page 2\n  ocloud compute instance list --json",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}

	// Add flags specific to the list command
	LimitFlag.Add(cmd)
	PageFlag.Add(cmd)
	JSONFlag.Add(cmd)

	return cmd
}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	// Get pagination parameters
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)

	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(logger.CmdLogger, 1, "Running instance list command in", "compartment", appCtx.CompartmentName, "limit", limit, "page", page, "json", useJSON)
	return instance.ListInstances(appCtx, limit, page, useJSON)
}
