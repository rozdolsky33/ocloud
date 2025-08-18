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
			return RunCreateCommand(cmd, appCtx)
		},
	}
	return cmd
}

// RunCreateCommand orchestrates the full flow. It calls TUI for selections,
// validates reachability, and spawns processes (SSH/tunnels) as needed.
func RunCreateCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	ctx := cmd.Context()

	svc, err := bastionSvc.NewService(appCtx)
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
	if b.ID == "" {
		return ErrAborted
	}
	tType, err := SelectTargetType(ctx, b.ID)
	if err != nil {
		return err
	}
	sType, err := SelectSessionType(ctx, b.ID)
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
