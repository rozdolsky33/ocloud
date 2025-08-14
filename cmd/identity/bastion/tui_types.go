// Package bastion Bubble Tea models and simple list UIs.
// Keep these UIs-only: no network calls or side effects here.
package bastion

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	instanceSvc "github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	okeSvc "github.com/rozdolsky33/ocloud/internal/services/compute/oke"
	autonomousdbSvc "github.com/rozdolsky33/ocloud/internal/services/database/autonomousdb"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
)

// SessionType identifies how the bastion session behaves.
type SessionType string

const (
	TypeManagedSSH     SessionType = "Managed SSH"
	TypePortForwarding SessionType = "Port-Forwarding"
)

// TargetType identifies what the session connects to.
type TargetType string

const (
	TargetOKE      TargetType = "OKE"
	TargetDatabase TargetType = "Database"
	TargetInstance TargetType = "Instance"
)

// BastionType identifies the top-level action.
type BastionType string

const (
	TypeBastion BastionType = "Bastion"
	TypeSession BastionType = "Session"
)

// ─── Type Selection ─────────────────────────────────────────────────────────────

type TypeSelectionModel struct {
	Cursor int
	Choice BastionType
	Types  []BastionType
}

func NewTypeSelectionModel() TypeSelectionModel {
	return TypeSelectionModel{
		Types:  []BastionType{TypeBastion, TypeSession},
		Cursor: 0,
	}
}
func (m TypeSelectionModel) Init() tea.Cmd { return nil }
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
			m.Cursor = (m.Cursor + 1) % len(m.Types)
		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Types) - 1
			}
		}
	}
	return m, nil
}
func (m TypeSelectionModel) View() string {
	var b strings.Builder
	b.WriteString("Select a type:\n\n")
	for i, t := range m.Types {
		mark := "( ) "
		if m.Cursor == i {
			mark = "(•) "
		}
		b.WriteString(mark + string(t) + "\n")
	}
	b.WriteString("\n(press q to quit)\n")
	return b.String()
}

// ─── Bastion Selection ──────────────────────────────────────────────────────────

type BastionModel struct {
	Cursor   int
	Choice   string
	Bastions []bastionSvc.Bastion
}

func NewBastionModel(bastions []bastionSvc.Bastion) BastionModel {
	return BastionModel{Bastions: bastions, Cursor: 0}
}
func (m BastionModel) Init() tea.Cmd { return nil }
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
			m.Cursor = (m.Cursor + 1) % len(m.Bastions)
		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Bastions) - 1
			}
		}
	}
	return m, nil
}
func (m BastionModel) View() string {
	var b strings.Builder
	b.WriteString("Select a bastion host:\n\n")
	for i, ba := range m.Bastions {
		mark := "( ) "
		if m.Cursor == i {
			mark = "(•) "
		}
		b.WriteString(mark + ba.Name + "\n")
	}
	b.WriteString("\n(press q to quit)\n")
	return b.String()
}

// ─── Session Type Selection ─────────────────────────────────────────────────────

type SessionTypeModel struct {
	Cursor    int
	Choice    SessionType
	Types     []SessionType
	BastionID string
}

func NewSessionTypeModel(bastionID string) SessionTypeModel {
	return SessionTypeModel{
		Types:     []SessionType{TypeManagedSSH, TypePortForwarding},
		Cursor:    0,
		BastionID: bastionID,
	}
}
func (m SessionTypeModel) Init() tea.Cmd { return nil }
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
			m.Cursor = (m.Cursor + 1) % len(m.Types)
		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Types) - 1
			}
		}
	}
	return m, nil
}
func (m SessionTypeModel) View() string {
	var b strings.Builder
	b.WriteString("Select a session type:\n\n")
	for i, t := range m.Types {
		mark := "( ) "
		if m.Cursor == i {
			mark = "(•) "
		}
		b.WriteString(mark + string(t) + "\n")
	}
	b.WriteString("\n(press q to quit)\n")
	return b.String()
}

// ─── Target Type Selection ──────────────────────────────────────────────────────

type TargetTypeModel struct {
	Cursor    int
	Choice    TargetType
	Types     []TargetType
	BastionID string
}

func NewTargetTypeModel(bastionID string, sessionType SessionType) TargetTypeModel {
	var types []TargetType
	if sessionType == TypeManagedSSH {
		types = []TargetType{TargetInstance}
	} else {
		types = []TargetType{TargetOKE, TargetDatabase, TargetInstance}
	}
	return TargetTypeModel{Types: types, Cursor: 0, BastionID: bastionID}
}
func (m TargetTypeModel) Init() tea.Cmd { return nil }
func (m TargetTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.Cursor = (m.Cursor + 1) % len(m.Types)
		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Types) - 1
			}
		}
	}
	return m, nil
}
func (m TargetTypeModel) View() string {
	var b strings.Builder
	b.WriteString("Select a target type:\n\n")
	for i, t := range m.Types {
		mark := "( ) "
		if m.Cursor == i {
			mark = "(•) "
		}
		b.WriteString(mark + string(t) + "\n")
	}
	b.WriteString("\n(press q to quit)\n")
	return b.String()
}

// ─── Fancy searchable lists (Instances / OKE / DB) ─────────────────────────────

type resourceItem struct {
	id, title, description string
}

func (i resourceItem) Title() string       { return i.title }
func (i resourceItem) Description() string { return i.description }
func (i resourceItem) FilterValue() string { return i.title + " " + i.description }

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

func NewInstanceListModelFancy(instances []instanceSvc.Instance) ResourceListModel {
	items := make([]list.Item, 0, len(instances))
	for _, inst := range instances {
		name := inst.Name
		if name == "" {
			if inst.Hostname != "" {
				name = inst.Hostname
			} else {
				name = inst.ID
			}
		}
		desc := inst.ID
		if inst.VcnName != "" {
			desc = inst.VcnName
			if inst.SubnetName != "" {
				desc += " · " + inst.SubnetName
			}
		}
		items = append(items, resourceItem{id: inst.ID, title: name, description: desc})
	}
	return newResourceList("Instances", items)
}

func NewOKEListModelFancy(clusters []okeSvc.Cluster) ResourceListModel {
	items := make([]list.Item, 0, len(clusters))
	for _, c := range clusters {
		desc := c.Version
		if c.PrivateEndpoint != "" {
			desc += " · PE"
		}
		items = append(items, resourceItem{id: c.ID, title: c.Name, description: desc})
	}
	return newResourceList("OKE Clusters", items)
}

func NewDBListModelFancy(dbs []autonomousdbSvc.AutonomousDatabase) ResourceListModel {
	items := make([]list.Item, 0, len(dbs))
	for _, d := range dbs {
		desc := d.PrivateEndpoint
		items = append(items, resourceItem{id: d.ID, title: d.Name, description: desc})
	}
	return newResourceList("Autonomous Databases", items)
}
