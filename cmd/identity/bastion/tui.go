// Package bastion provides commands for managing OCI Bastions.
// This file contains the Terminal User Interface (TUI) implementation
// using the Bubble Tea framework for interactive selection.
package bastion

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
)

// SessionType represents the type of bastion session
type SessionType string

const (
	// TypeManagedSSH represents a managed SSH session
	TypeManagedSSH SessionType = "Managed SSH"

	// TypePortForwarding represents a port-forwarding session
	TypePortForwarding SessionType = "Port-Forwarding"
)

// TargetType represents the target for port forwarding
type TargetType string

const (
	// TargetOKE represents an OKE target for port forwarding
	TargetOKE TargetType = "OKE"

	// TargetDatabase represents a database target for port forwarding
	TargetDatabase TargetType = "Database"

	// TargetInstance represents an instance target for port forwarding.
	TargetInstance TargetType = "Instance"
)

// BastionType represents the type of bastion operation
// This is used to differentiate between direct bastion creation and session-based creation
type BastionType string

const (
	// TypeBastion represents a direct bastion creation
	// This is the standard mode for creating bastions
	TypeBastion BastionType = "Bastion"

	// TypeSession represents a session-based bastion creation
	// This mode uses the dummy bastion list for selection
	TypeSession BastionType = "Session"
)

// TypeSelectionModel represents the TUI model for selecting between Bastion and Session types
// It implements the tea.Model interface required by Bubble Tea
type TypeSelectionModel struct {
	Cursor int
	Choice BastionType
	Types  []BastionType
}

// Init initializes the model
// This is part of the tea.Model interface
func (m TypeSelectionModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model accordingly
// This is part of the tea.Model interface
func (m TypeSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			if m.Cursor >= 0 && m.Cursor < len(m.Types) {
				m.Choice = m.Types[m.Cursor]
			}
			return m, tea.Quit

		case "down", "j":
			// Move the cursor down (with wrap-around)
			m.Cursor++
			if m.Cursor >= len(m.Types) {
				m.Cursor = 0
			}

		case "up", "k":
			// Move the cursor up (with wrap-around)
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Types) - 1
			}
		}
	}

	return m, nil
}

// View renders the current UI state as a string
func (m TypeSelectionModel) View() string {
	s := strings.Builder{}
	s.WriteString("Select a type:\n\n")

	for i, t := range m.Types {
		if m.Cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(string(t))
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

// NewTypeSelectionModel creates a new TypeSelectionModel
func NewTypeSelectionModel() TypeSelectionModel {
	return TypeSelectionModel{
		Types:  []BastionType{TypeBastion, TypeSession},
		Cursor: 0,
	}
}

// BastionModel represents the TUI model for bastion selection
// It implements the tea.Model interface required by Bubble Tea
type BastionModel struct {
	Cursor   int
	Choice   string
	Bastions []bastionSvc.Bastion
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
			return m, tea.Quit

		case "enter":
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
	s.WriteString("Select a bastion host to create a session:\n\n")

	for i, bastion := range m.Bastions {
		if m.Cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
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

// SessionTypeModel represents the TUI model for session type selection
// It implements the tea.Model interface required by Bubble Tea
type SessionTypeModel struct {
	Cursor    int
	Choice    SessionType
	Types     []SessionType
	BastionID string
}

// Init initializes the model
// This is part of the tea.Model interface
func (m SessionTypeModel) Init() tea.Cmd {
	// No initialization needed
	return nil
}

// Update handles messages and updates the model accordingly
// This is part of the tea.Model interface
func (m SessionTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			if m.Cursor >= 0 && m.Cursor < len(m.Types) {
				m.Choice = m.Types[m.Cursor]
			}
			return m, tea.Quit

		case "down", "j":
			// Move the cursor down (with wrap-around)
			m.Cursor++
			if m.Cursor >= len(m.Types) {
				m.Cursor = 0
			}

		case "up", "k":
			// Move the cursor up (with wrap-around)
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Types) - 1
			}
		}
	}

	return m, nil
}

// View renders the current UI state as a string
// This is part of the tea.Model interface
func (m SessionTypeModel) View() string {
	s := strings.Builder{}
	s.WriteString("Select a session type:\n\n")

	for i, t := range m.Types {
		if m.Cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(string(t))
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

// NewSessionTypeModel creates a new SessionTypeModel with the provided bastion ID
func NewSessionTypeModel(bastionID string) SessionTypeModel {
	return SessionTypeModel{
		Types:     []SessionType{TypeManagedSSH, TypePortForwarding},
		Cursor:    0,
		BastionID: bastionID,
	}
}

// TargetTypeModel represents the TUI model for target type selection
// It implements the tea.Model interface required by Bubble Tea
type TargetTypeModel struct {
	Cursor    int
	Choice    TargetType
	Types     []TargetType
	BastionID string
}

// Init initializes the model
// This is part of the tea.Model interface
func (m TargetTypeModel) Init() tea.Cmd {
	// No initialization needed
	return nil
}

// Update handles messages and updates the model accordingly
// This is part of the tea.Model interface
func (m TargetTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Set the choice and exit
			if m.Cursor >= 0 && m.Cursor < len(m.Types) {
				m.Choice = m.Types[m.Cursor]
			}
			return m, tea.Quit

		case "down", "j":
			// Move the cursor down (with wrap-around)
			m.Cursor++
			if m.Cursor >= len(m.Types) {
				m.Cursor = 0
			}

		case "up", "k":
			// Move the cursor up (with wrap-around)
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Types) - 1
			}
		}
	}

	return m, nil
}

// View renders the current UI state as a string
// This is part of the tea.Model interface
func (m TargetTypeModel) View() string {
	s := strings.Builder{}
	s.WriteString("Select a target type:\n\n")

	// Render each target type option with a cursor indicator
	for i, t := range m.Types {
		if m.Cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(string(t))
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

// NewTargetTypeModel creates a new TargetTypeModel with the provided bastion ID
func NewTargetTypeModel(bastionID string) TargetTypeModel {
	return TargetTypeModel{
		Types:     []TargetType{TargetOKE, TargetDatabase, TargetInstance},
		Cursor:    0,
		BastionID: bastionID,
	}
}
