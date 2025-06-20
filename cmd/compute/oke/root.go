package oke

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
)

// NewOKECmd creates a new command for instance-related operations
func NewOKECmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "oke",
		Short:         "Manage OCI Kubernetes Engine (OKE)",
		Long:          "Manage Oracle Cloud Infrastructure Kubernetes Engine - list all or find clusters by name pattern.",
		Example:       "  ocloud compute oke list\n  ocloud compute oke find myoke",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd(appCtx))
	cmd.AddCommand(NewFindCmd(appCtx))

	return cmd
}
