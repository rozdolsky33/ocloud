package vcn

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	cfgflags "github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	netvcn "github.com/rozdolsky33/ocloud/internal/services/network/vcn"
	"github.com/spf13/cobra"
)

func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "Lists VCNs in a compartment",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cmd, appCtx)
		},
	}
	cmd.Flags().Bool("gateways", false, "Display gateways")
	cmd.Flags().Bool("subnets", false, "Display subnets")
	cmd.Flags().Bool("nsg", false, "Display network security groups")
	cmd.Flags().Bool("route", false, "Display route tables")
	cmd.Flags().Bool("security-list", false, "Display security lists")
	return cmd
}

func runListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	useJSON := cfgflags.GetBoolFlag(cmd, cfgflags.FlagNameJSON, false)

	gateways, _ := cmd.Flags().GetBool("gateways")
	subnets, _ := cmd.Flags().GetBool("subnets")
	nsgs, _ := cmd.Flags().GetBool("nsg")
	routes, _ := cmd.Flags().GetBool("route")
	securityLists, _ := cmd.Flags().GetBool("security-list")

	logger.LogWithLevel(logger.CmdLogger, logger.Debug, "Running network vcn list", "json", useJSON)

	return netvcn.ListVCNs(appCtx, useJSON, gateways, subnets, nsgs, routes, securityLists)
}
