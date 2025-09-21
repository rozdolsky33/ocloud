package tui

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

// ErrCancelled is returned when the user quits without confirming.
var ErrCancelled = errors.New("selection cancelled")

// Run returns the selected ID, or ErrCancelled if the user quit without confirming.
func Run(m Model) (string, error) {
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}
	if mm, ok := finalModel.(Model); ok {
		if !mm.confirmed || mm.choice == "" {
			return "", ErrCancelled
		}
		return mm.choice, nil
	}
	return "", ErrCancelled
}
