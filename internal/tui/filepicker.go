package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// FileInfo represents a file or directory entry.
type FileInfo struct {
	Name  string
	Path  string
	Size  int64
	IsDir bool
}

// fileItem implements list.Item for the file picker.
type fileItem struct {
	info FileInfo
}

func (i fileItem) Title() string {
	if i.info.IsDir {
		return i.info.Name + "/"
	}
	return i.info.Name
}

func (i fileItem) Description() string {
	if i.info.IsDir {
		return "Directory"
	}
	return humanizeBytes(i.info.Size)
}

func (i fileItem) FilterValue() string { return i.info.Name }

// FilePickerModel is a TUI for selecting a file from the current directory.
type FilePickerModel struct {
	list       list.Model
	currentDir string
	choice     string
	confirmed  bool
	keys       KeyMap
	showHidden bool
}

// NewFilePickerModel creates a file picker starting from the given directory.
func NewFilePickerModel(startDir string) (FilePickerModel, error) {
	absDir, err := filepath.Abs(startDir)
	if err != nil {
		return FilePickerModel{}, err
	}

	m := FilePickerModel{
		currentDir: absDir,
		keys:       DefaultKeyMap(),
		showHidden: false,
	}

	items, err := m.readDir()
	if err != nil {
		return FilePickerModel{}, err
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = fmt.Sprintf("Select file from: %s", absDir)
	l.SetShowTitle(true)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)

	m.list = l
	return m, nil
}

func (m *FilePickerModel) readDir() ([]list.Item, error) {
	entries, err := os.ReadDir(m.currentDir)
	if err != nil {
		return nil, err
	}

	var files []FileInfo

	// Add parent directory if not at root
	if m.currentDir != "/" {
		files = append(files, FileInfo{
			Name:  "..",
			Path:  filepath.Dir(m.currentDir),
			IsDir: true,
		})
	}

	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden files unless showHidden is enabled
		if !m.showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			Name:  name,
			Path:  filepath.Join(m.currentDir, name),
			Size:  info.Size(),
			IsDir: entry.IsDir(),
		})
	}

	// Sort: directories first, then files, both alphabetically
	sort.Slice(files, func(i, j int) bool {
		// Keep ".." at the top
		if files[i].Name == ".." {
			return true
		}
		if files[j].Name == ".." {
			return false
		}
		// Directories before files
		if files[i].IsDir && !files[j].IsDir {
			return true
		}
		if !files[i].IsDir && files[j].IsDir {
			return false
		}
		// Alphabetically within same type
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	items := make([]list.Item, len(files))
	for i, f := range files {
		items[i] = fileItem{info: f}
	}

	return items, nil
}

func (m FilePickerModel) Init() tea.Cmd { return nil }

func (m FilePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			m.confirmed = false
			return m, tea.Quit
		}
		if key.Matches(msg, m.keys.Confirm) {
			if item, ok := m.list.SelectedItem().(fileItem); ok {
				if item.info.IsDir {
					// Navigate into directory
					m.currentDir = item.info.Path
					items, err := m.readDir()
					if err == nil {
						m.list.SetItems(items)
						m.list.Title = fmt.Sprintf("Select file from: %s", m.currentDir)
						m.list.ResetSelected()
					}
					return m, nil
				}
				// File selected
				m.choice = item.info.Path
				m.confirmed = true
				return m, tea.Quit
			}
		}
		// Toggle hidden files with '.'
		if msg.String() == "." {
			m.showHidden = !m.showHidden
			items, err := m.readDir()
			if err == nil {
				m.list.SetItems(items)
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m FilePickerModel) View() string   { return m.list.View() }
func (m FilePickerModel) Choice() string { return m.choice }

// RunFilePicker runs the file picker and returns the selected file path.
func RunFilePicker(m FilePickerModel) (string, error) {
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}
	if mm, ok := finalModel.(FilePickerModel); ok {
		if !mm.confirmed || mm.choice == "" {
			return "", ErrCancelled
		}
		return mm.choice, nil
	}
	return "", ErrCancelled
}

// humanizeBytes formats bytes in human readable format.
func humanizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
