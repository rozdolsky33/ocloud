package image

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// resourceItem defines a resource item for a list.
type resourceItem struct {
	id, title, description string
}

func (i resourceItem) Title() string       { return i.title }
func (i resourceItem) Description() string { return i.description }
func (i resourceItem) FilterValue() string { return i.title + " " + i.description }

// ResourceListModel defines a TUI model for displaying a list of resources.
type ResourceListModel struct {
	list   list.Model
	choice string
	keys   struct {
		confirm key.Binding
		quit    key.Binding
	}
}

func (m ResourceListModel) Init() tea.Cmd { return nil }
func (m ResourceListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.quit) {
			return m, tea.Quit
		}
		if key.Matches(msg, m.keys.confirm) {
			if it, ok := m.list.SelectedItem().(resourceItem); ok {
				m.choice = it.id
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ResourceListModel) View() string   { return m.list.View() }
func (m ResourceListModel) Choice() string { return m.choice }

func newResourceList(title string, items []list.Item) ResourceListModel {
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = title
	l.SetShowTitle(true)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)

	rm := ResourceListModel{list: l}
	rm.keys.confirm = key.NewBinding(key.WithKeys("enter"))
	rm.keys.quit = key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"))
	return rm
}

// NewImageListModelFancy creates a ResourceListModel to display instances in a searchable and interactive list.
// It transforms each instance into a resourceItem with its name, ID, and VCN name as attributes.
func NewImageListModelFancy(images []Image) ResourceListModel {
	items := make([]list.Item, 0, len(images))
	for _, img := range images {
		name := img.Name
		desc := img.OperatingSystem
		id := img.ID
		items = append(items, resourceItem{id: id, title: name, description: desc})
	}
	return newResourceList("Images", items)
}
