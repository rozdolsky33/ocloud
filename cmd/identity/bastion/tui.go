// Package bastion provides commands for managing OCI Bastions.
// This file contains the Terminal User Interface (TUI) implementation
// using the Bubble Tea framework for interactive selection.
package bastion

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
)

// BastionModel represents the TUI model for bastion selection
// It implements the tea.Model interface required by Bubble Tea
type BastionModel struct {
	Cursor   int                  // Current cursor position
	Choice   string               // Selected bastion ID
	Bastions []bastionSvc.Bastion // List of available bastions
}

// Init initializes the model
// This is part of the tea.Model interface
func (m BastionModel) Init() tea.Cmd {
	// No initialization needed
	return nil
}

// Update handles messages and updates the model accordingly
// This is part of the tea.Model interface
func (m BastionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			// Exit without selection
			return m, tea.Quit

		case "enter":
			// Set the choice and exit
			if m.Cursor >= 0 && m.Cursor < len(m.Bastions) {
				m.Choice = m.Bastions[m.Cursor].ID
			}
			return m, tea.Quit

		case "down", "j":
			// Move the cursor down (with wrap-around)
			m.Cursor++
			if m.Cursor >= len(m.Bastions) {
				m.Cursor = 0
			}

		case "up", "k":
			// Move the cursor up (with wrap-around)
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Bastions) - 1
			}
		}
	}

	return m, nil
}

// View renders the current UI state as a string
// This is part of the tea.Model interface
func (m BastionModel) View() string {
	s := strings.Builder{}
	s.WriteString("Select a Bastion to create:\n\n")

	// Render each bastion option with a cursor indicator
	for i, bastion := range m.Bastions {
		if m.Cursor == i {
			s.WriteString("(•) ") // Selected item
		} else {
			s.WriteString("( ) ") // Unselected item
		}
		s.WriteString(bastion.Name)
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

// NewBastionModel creates a new BastionModel with the provided bastions
func NewBastionModel(bastions []bastionSvc.Bastion) BastionModel {
	return BastionModel{
		Bastions: bastions,
		Cursor:   0,
	}
}
