package auth

import (
	"fmt"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// BrowserOption represents a selectable browser option.
type BrowserOption struct {
	Label string
	Value string
}

// BrowserSelectionModel is a minimal Bubble Tea model for selecting a browser.
type BrowserSelectionModel struct {
	Cursor  int
	Choice  int
	Options []BrowserOption
}

func newBrowserSelectionModel() BrowserSelectionModel {

	isMac := runtime.GOOS == "darwin"

	// Helper to choose a value per OS
	firefoxVal := "firefox"
	chromeVal := "google-chrome"
	braveVal := "brave"

	safariVal := "open -b com.apple.Safari"
	chromiumVal := "chromium"
	if isMac {
		firefoxVal = "open -b org.mozilla.firefox"
		chromeVal = "open -b com.google.Chrome"
		braveVal = "open -b com.brave.Browser"
		chromiumVal = "open -b org.chromium.Chromium"
	}

	opts := []BrowserOption{
		{Label: "Firefox", Value: firefoxVal},
		{Label: "Google Chrome", Value: chromeVal},
		{Label: "Chromium", Value: chromiumVal},
		{Label: "Brave", Value: braveVal},
	}
	// Safari only listed as an explicit option on macOS
	if isMac {
		opts = append(opts, BrowserOption{Label: "Safari", Value: safariVal})
	}

	return BrowserSelectionModel{Cursor: 0, Choice: -1, Options: opts}
}

func (m BrowserSelectionModel) Init() tea.Cmd { return nil }

func (m BrowserSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			if m.Cursor >= 0 && m.Cursor < len(m.Options) {
				m.Choice = m.Cursor
			}
			return m, tea.Quit
		case "down", "j":
			m.Cursor = (m.Cursor + 1) % len(m.Options)
		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Options) - 1
			}
		}
	}
	return m, nil
}

func (m BrowserSelectionModel) View() string {
	var b strings.Builder
	b.WriteString("Select a browser to use for OCI authentication:\n\n")
	for i, opt := range m.Options {
		mark := "( ) "
		if m.Cursor == i {
			mark = "(•) "
		}
		b.WriteString(mark + opt.Label + "\n")
	}
	b.WriteString("\nenter: confirm    q: quit\n")
	return b.String()
}

// RunBrowserPicker runs the TUI and returns the browser value to set in BROWSER and a boolean indicating whether to set it.
func RunBrowserPicker() (value string, set bool, err error) {
	model := newBrowserSelectionModel()
	p := tea.NewProgram(model)
	res, err := p.StartReturningModel()
	if err != nil {
		return "", false, err
	}
	m, ok := res.(BrowserSelectionModel)
	if !ok {
		return "", false, fmt.Errorf("unexpected model type")
	}
	if m.Choice < 0 || m.Choice >= len(m.Options) {
		return "", false, nil
	}
	chosen := m.Options[m.Choice]
	switch chosen.Value {
	case "__KEEP__":
		return "", false, nil
	case "__UNSET__":
		return "__UNSET__", true, nil
	default:
		return chosen.Value, true, nil
	}
}
