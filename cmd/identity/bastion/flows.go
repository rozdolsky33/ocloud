// Package bastion Flows (orchestrators) that stitch together services, TUI, and side effects.
// These are thin and testable: they take ctx and collaborators; no globals.
package bastion

import (
	"context"
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
)

// SelectBastionType runs a simple TUI to choose between Bastion mgmt or Session.
func SelectBastionType(ctx context.Context) (BastionType, error) {
	m := NewTypeSelectionModel()
	p := tea.NewProgram(m, tea.WithContext(ctx))
	res, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("type selection TUI: %w", err)
	}
	out, ok := res.(TypeSelectionModel)
	if !ok || out.Choice == "" {
		return "", ErrAborted
	}
	return out.Choice, nil
}

// SelectBastion lists and filters ACTIVE bastions, then runs a picker TUI.
func SelectBastion(ctx context.Context, svc *bastionSvc.Service, t BastionType) (bastionSvc.Bastion, error) {
	if t != TypeSession {
		return bastionSvc.Bastion{}, nil
	}
	list, err := svc.List(ctx)
	if err != nil {
		return bastionSvc.Bastion{}, fmt.Errorf("list bastions: %w", err)
	}
	list = slices.DeleteFunc(list, func(b bastionSvc.Bastion) bool {
		return b.LifecycleState != "ACTIVE"
	})

	m := NewBastionModel(list)
	p := tea.NewProgram(m, tea.WithContext(ctx))
	res, err := p.Run()
	if err != nil {
		return bastionSvc.Bastion{}, fmt.Errorf("bastion selection TUI: %w", err)
	}
	out, ok := res.(BastionModel)
	if !ok || out.Choice == "" {
		return bastionSvc.Bastion{}, ErrAborted
	}
	for _, b := range list {
		if b.OCID == out.Choice {
			return b, nil
		}
	}
	return bastionSvc.Bastion{}, fmt.Errorf("selected bastion not found")
}

// SelectTargetType provides a TUI to select a target type associated with the given bastion ID and returns the selection.
func SelectTargetType(ctx context.Context, bastionID string) (TargetType, error) {
	m := NewTargetTypeModel(bastionID)
	p := tea.NewProgram(m, tea.WithContext(ctx))
	res, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("target type TUI: %w", err)
	}
	out, ok := res.(TargetTypeModel)
	if !ok || out.Choice == "" {
		return "", ErrAborted
	}
	return out.Choice, nil
}

// SelectSessionType chooses a session type for the selected bastion.
func SelectSessionType(ctx context.Context, bastionID string) (SessionType, error) {
	m := NewSessionTypeModel(bastionID)
	p := tea.NewProgram(m, tea.WithContext(ctx))
	res, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("session type TUI: %w", err)
	}
	out, ok := res.(SessionTypeModel)
	if !ok || out.Choice == "" {
		return "", ErrAborted
	}
	return out.Choice, nil
}

// ConnectTarget switches to the correct flow for the chosen target.
func ConnectTarget(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
	b bastionSvc.Bastion, sType SessionType, tType TargetType) error {

	switch tType {
	case TargetInstance:
		return connectInstance(ctx, appCtx, svc, b, sType)
	case TargetDatabase:
		return connectDatabase(ctx, appCtx, svc, b, sType)
	case TargetOKE:
		return connectOKE(ctx, appCtx, svc, b, sType)
	default:
		fmt.Printf("Prepared %s session on %s (%s) -> %s\n", sType, b.DisplayName, b.OCID, tType)
		return nil
	}
}
