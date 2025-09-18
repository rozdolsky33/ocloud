package vcn

import (
	vcnFlags "github.com/rozdolsky33/ocloud/cmd/shared/flags"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	netvcn "github.com/rozdolsky33/ocloud/internal/services/network/vcn"
	"github.com/spf13/cobra"
)

// Long description for the get command
var getLong = `
Display a VCN summary by OCID.

This command fetches a Virtual Cloud Network (VCN) by its OCID and prints a concise summary
including name, state, compartment, CIDR blocks, DNS label/domain, DHCP options, creation time,
and tags. Use --json to get the raw JSON output instead of the formatted summary.`

// Examples for the get command
var getExamples = `
 `

// NewGetCmd returns "vcn get" command.
func NewGetCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Get VCN summary by OCID",
		Long:          getLong,
		Example:       getExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGetCommand(cmd, appCtx)
		},
	}
	vcnFlags.LimitFlag.Add(cmd)
	vcnFlags.PageFlag.Add(cmd)

	return cmd
}

// RunGetCommand executes the get logic
func RunGetCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	limit := flags.GetIntFlag(cmd, flags.FlagNameLimit, vcnFlags.FlagDefaultLimit)
	page := flags.GetIntFlag(cmd, flags.FlagNamePage, vcnFlags.FlagDefaultPage)
	useJSON := flags.GetBoolFlag(cmd, flags.FlagNameJSON, false)

	return netvcn.GetVCNs(appCtx, limit, page, useJSON)
}
