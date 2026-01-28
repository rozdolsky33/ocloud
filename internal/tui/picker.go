package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PickerOption represents an option in the picker.
type PickerOption struct {
	ID          string
	Label       string
	Description string
}

// PickerModel is a radio button style picker for selecting a single option.
type PickerModel struct {
	title     string
	options   []PickerOption
	cursor    int
	choice    string
	confirmed bool
	keys      KeyMap
}

// NewPickerModel creates a new picker model with the given title and options.
func NewPickerModel(title string, options []PickerOption) PickerModel {
	return PickerModel{
		title:   title,
		options: options,
		cursor:  0,
		keys:    DefaultKeyMap(),
	}
}

func (m PickerModel) Init() tea.Cmd { return nil }

func (m PickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			m.confirmed = false
			return m, tea.Quit
		}
		if key.Matches(msg, m.keys.Confirm) {
			if len(m.options) > 0 {
				m.choice = m.options[m.cursor].ID
				m.confirmed = true
			}
			return m, tea.Quit
		}
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m PickerModel) View() string {
	var b strings.Builder

	// Styles
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	unselectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).MarginLeft(4)

	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n\n")

	for i, opt := range m.options {
		var radio string
		var labelStyle lipgloss.Style
		if i == m.cursor {
			radio = "(•)"
			labelStyle = selectedStyle
		} else {
			radio = "( )"
			labelStyle = unselectedStyle
		}

		b.WriteString(fmt.Sprintf("%s %s\n", radio, labelStyle.Render(opt.Label)))
		if opt.Description != "" {
			b.WriteString(descStyle.Render(opt.Description))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("↑/↓: navigate • enter: select • esc: cancel"))

	return b.String()
}

func (m PickerModel) Choice() string { return m.choice }

// RunPicker runs the picker and returns the selected option ID.
func RunPicker(m PickerModel) (string, error) {
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}
	if mm, ok := finalModel.(PickerModel); ok {
		if !mm.confirmed || mm.choice == "" {
			return "", ErrCancelled
		}
		return mm.choice, nil
	}
	return "", ErrCancelled
}
