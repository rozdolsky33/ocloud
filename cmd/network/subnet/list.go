package subnet

import (
	paginationFlags "github.com/rozdolsky33/ocloud/cmd/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/network/subnet"
	"github.com/spf13/cobra"
)

// NewListCmd creates a new "list" command for listing all subnets in a specified tenancy or compartment.
// Accepts an ApplicationContext for accessing configuration and dependencies.
// Adds pagination flags for controlling the number of results returned and the page to retrieve.
// Executes the RunListCommand function when invoked.
func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all Subnets in the specified tenancy or compartment",
		Long:          "",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}
	// Add flags specific to the list command
	paginationFlags.LimitFlag.Add(cmd)
	paginationFlags.PageFlag.Add(cmd)

	return cmd

}

// RunListCommand handles the execution of the list command
func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	// Get pagination parameters
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, paginationFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, paginationFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)

	// Use LogWithLevel to ensure debug logs work with shorthand flags
	logger.LogWithLevel(logger.CmdLogger, 1, "Running subnet list command in", "compartment", appCtx.CompartmentName, "json", useJSON)
	return subnet.ListSubnets(appCtx, useJSON, limit, page)
}
