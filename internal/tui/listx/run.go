package listx

import tea "github.com/charmbracelet/bubbletea"

// Run returns the selected ID (or empty if none).
func Run(m Model) (string, error) {
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}
	if mm, ok := finalModel.(Model); ok {
		return mm.Choice(), nil
	}
	return "", nil
}
