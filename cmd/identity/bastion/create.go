// Package bastion Command wiring and orchestration for "bastion creates".
package bastion

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/app"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/spf13/cobra"
)

// NewCreateCmd returns "bastion create".
func NewCreateCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "create",
		Aliases:       []string{"c"},
		Short:         "Create a Bastion or a Session",
		Long:          "Interactively create a session on a selected bastion and target (Instance, OKE, Database).",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCreateCommand(cmd, appCtx)
		},
	}
	return cmd
}

// runCreateCommand orchestrates the full flow. It calls TUI for selections.
func runCreateCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	ctx := cmd.Context()

	svc, err := bastionSvc.NewServiceFromAppContext(appCtx)
	if err != nil {
		return fmt.Errorf("create bastion service: %w", err)
	}

	choice, err := SelectBastionType(ctx)
	if err != nil {
		return err
	}
	if choice == "" {
		return ErrAborted
	}
	if choice == TypeBastion {
		util.ShowConstructionAnimation()
		return nil
	}
	b, err := SelectBastion(ctx, svc, choice)
	if err != nil {
		return err
	}
	if b.OCID == "" {
		return ErrAborted
	}
	tType, err := SelectTargetType(ctx, b.OCID)
	if err != nil {
		return err
	}
	sType, err := SelectSessionType(ctx, b.OCID)
	if err != nil {
		return err
	}
	if sType == "" {
		return ErrAborted
	}
	if tType == "" {
		return ErrAborted
	}

	return ConnectTarget(ctx, appCtx, svc, b, sType, tType)
}
