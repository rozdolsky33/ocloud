package listx

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// ResourceItemData is a lightweight DTO for the list.
type ResourceItemData struct {
	ID, Title, Description string
}

// resourceItem implements bubbles/list.Item.
type resourceItem struct {
	id, title, description string
}

func (i resourceItem) Title() string       { return i.title }
func (i resourceItem) Description() string { return i.description }
func (i resourceItem) FilterValue() string { return i.title + " " + i.description }

// KeyMap defines key bindings (export if you want callers to override).
type KeyMap struct {
	Confirm key.Binding
	Quit    key.Binding
}

// DefaultKeyMap returns a sensible default.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Confirm: key.NewBinding(key.WithKeys("enter")),
		Quit:    key.NewBinding(key.WithKeys("q", "esc", "ctrl+c")),
	}
}

// Model is a reusable Bubble Tea model for a searchable list of resources.
type Model struct {
	list   list.Model
	choice string
	keys   KeyMap
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}
		if key.Matches(msg, m.keys.Confirm) {
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

func (m Model) View() string   { return m.list.View() }
func (m Model) Choice() string { return m.choice }

// NewModel creates a list model from arbitrary data by using an adapter.
//
// T is your domain type (Image, Instance, Node, etc.).
// adapter maps T -> ResourceItemData (id/title/description).
func NewModel[T any](title string, data []T, adapter func(T) ResourceItemData) Model {
	items := make([]list.Item, 0, len(data))
	for _, d := range data {
		ri := adapter(d)
		items = append(items, resourceItem{id: ri.ID, title: ri.Title, description: ri.Description})
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = title
	l.SetShowTitle(true)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)

	m := Model{list: l, keys: DefaultKeyMap()}
	return m
}
