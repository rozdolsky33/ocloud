package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TextInputModel is a Bubble Tea model for a single-line text input with a prompt.
type TextInputModel struct {
	textInput textinput.Model
	header    string
	footer    string
	confirmed bool
	quitting  bool
}

// NewTextInputModel creates a text input TUI with the given prompt, placeholder, and default value.
func NewTextInputModel(header, placeholder, defaultValue string) TextInputModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 512
	ti.Width = 60
	if defaultValue != "" {
		ti.SetValue(defaultValue)
	}

	return TextInputModel{
		textInput: ti,
		header:    header,
		footer:    "(enter to confirm, esc to cancel)",
	}
}

func (m TextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.confirmed = true
			m.quitting = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m TextInputModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170"))
	footerStyle := lipgloss.NewStyle().Faint(true)

	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		headerStyle.Render(m.header),
		"",
		m.textInput.View(),
		"",
		footerStyle.Render(m.footer),
		"",
	)
}

// Value returns the text input value.
func (m TextInputModel) Value() string { return m.textInput.Value() }

// RunTextInput runs the text input TUI and returns the entered value.
// Returns ErrCancelled if the user quit without confirming.
func RunTextInput(m TextInputModel) (string, error) {
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}
	if mm, ok := finalModel.(TextInputModel); ok {
		if !mm.confirmed {
			return "", ErrCancelled
		}
		return mm.Value(), nil
	}
	return "", ErrCancelled
}
