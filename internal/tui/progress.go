package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressMsg is sent to update progress.
type ProgressMsg float64

// ProgressErrMsg is sent when an error occurs.
type ProgressErrMsg struct{ Err error }

// ProgressDoneMsg is sent when the operation completes.
type ProgressDoneMsg struct{}

// ProgressModel is a TUI model for showing progress of an operation.
type ProgressModel struct {
	progress  progress.Model
	title     string
	status    string
	percent   float64
	err       error
	done      bool
	bytesInfo string // e.g., "10.5 MiB / 50 MiB"
	extraInfo string // e.g., "Part 2/5"
}

// NewProgressModel creates a new progress model with the given title.
func NewProgressModel(title string) ProgressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)
	return ProgressModel{
		progress: p,
		title:    title,
		status:   "Starting...",
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return nil
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Allow quitting with ctrl+c
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 10
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}
		return m, nil

	case ProgressMsg:
		m.percent = float64(msg)
		if m.percent >= 1.0 {
			m.done = true
			m.status = "Complete!"
			return m, tea.Quit
		}
		return m, nil

	case ProgressErrMsg:
		m.err = msg.Err
		m.status = fmt.Sprintf("Error: %v", msg.Err)
		return m, tea.Quit

	case ProgressDoneMsg:
		m.done = true
		m.percent = 1.0
		m.status = "Complete!"
		return m, tea.Quit

	case ProgressUpdateMsg:
		m.percent = msg.Percent
		m.bytesInfo = msg.BytesInfo
		m.extraInfo = msg.ExtraInfo
		m.status = msg.Status
		if m.percent >= 1.0 {
			m.done = true
			return m, tea.Quit
		}
		return m, nil
	}

	return m, nil
}

func (m ProgressModel) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n\n")

	b.WriteString(m.progress.ViewAs(m.percent))
	b.WriteString("\n\n")

	// Show bytes info if available
	if m.bytesInfo != "" {
		b.WriteString(infoStyle.Render(m.bytesInfo))
		if m.extraInfo != "" {
			b.WriteString(" • ")
			b.WriteString(infoStyle.Render(m.extraInfo))
		}
		b.WriteString("\n")
	}

	b.WriteString(statusStyle.Render(m.status))
	b.WriteString("\n")

	return b.String()
}

// ProgressUpdateMsg contains detailed progress information.
type ProgressUpdateMsg struct {
	Percent   float64
	BytesInfo string
	ExtraInfo string
	Status    string
}

// ProgressRunner runs an operation with a progress TUI.
type ProgressRunner struct {
	program *tea.Program
	model   ProgressModel
}

// NewProgressRunner creates a new progress runner with the given title.
func NewProgressRunner(title string) *ProgressRunner {
	model := NewProgressModel(title)
	return &ProgressRunner{
		model: model,
	}
}

// Start starts the progress TUI. Call this before starting the operation.
func (r *ProgressRunner) Start() {
	r.program = tea.NewProgram(r.model)
}

// Run runs the progress TUI and blocks until complete.
func (r *ProgressRunner) Run() error {
	if r.program == nil {
		r.Start()
	}
	finalModel, err := r.program.Run()
	if err != nil {
		return err
	}
	if m, ok := finalModel.(ProgressModel); ok && m.err != nil {
		return m.err
	}
	return nil
}

// UpdateProgress sends a progress update to the TUI.
func (r *ProgressRunner) UpdateProgress(percent float64, bytesInfo, extraInfo, status string) {
	if r.program != nil {
		r.program.Send(ProgressUpdateMsg{
			Percent:   percent,
			BytesInfo: bytesInfo,
			ExtraInfo: extraInfo,
			Status:    status,
		})
	}
}

// SendError sends an error to the TUI.
func (r *ProgressRunner) SendError(err error) {
	if r.program != nil {
		r.program.Send(ProgressErrMsg{Err: err})
	}
}

// SendDone signals that the operation is complete.
func (r *ProgressRunner) SendDone() {
	if r.program != nil {
		r.program.Send(ProgressDoneMsg{})
	}
}

// Program returns the underlying tea.Program for sending custom messages.
func (r *ProgressRunner) Program() *tea.Program {
	return r.program
}
