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
		Short:         "Create a Bastion session",
		Long:          "Interactively create a session on a selected bastion and target (instance, OKE, database).",
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

	choice, err := SelectBastionType(ctx) // tui (ctx-aware)
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

	b, err := SelectBastion(ctx, svc, choice) // flow: lists active, TUI picker
	if err != nil {
		return err
	}
	if b.ID == "" {
		return ErrAborted
	}

	sType, err := SelectSessionType(ctx, b.ID) // tui
	if err != nil {
		return err
	}
	if sType == "" {
		return ErrAborted
	}

	tType, err := SelectTargetType(ctx, b.ID, sType) // tui
	if err != nil {
		return err
	}
	if tType == "" {
		return ErrAborted
	}

	return ConnectTarget(ctx, appCtx, svc, b, sType, tType)
}
