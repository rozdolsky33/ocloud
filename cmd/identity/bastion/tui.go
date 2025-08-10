// Package bastion provides commands for managing OCI Bastions.
// This file contains the Terminal User Interface (TUI) implementation
// using the Bubble Tea framework for interactive selection.
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

// SessionType represents the type of bastion session
type SessionType string

const (
	TypeManagedSSH     SessionType = "Managed SSH"
	TypePortForwarding SessionType = "Port-Forwarding"
)

// TargetType represents the target for port forwarding
type TargetType string

const (
	TargetOKE      TargetType = "OKE"
	TargetDatabase TargetType = "Database"
	TargetInstance TargetType = "Instance"
)

// BastionType represents the type of bastion operation
type BastionType string

const (
	TypeBastion BastionType = "Bastion"
	TypeSession BastionType = "Session"
)

// TypeSelectionModel represents the TUI model for selecting between Bastion and Session types
// It implements the tea.Model interface required by Bubble Tea
type TypeSelectionModel struct {
	Cursor int
	Choice BastionType
	Types  []BastionType
}

// Init initializes the model
// This is part of the tea.Model interface
func (m TypeSelectionModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model accordingly
// This is part of the tea.Model interface
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
			m.Cursor++
			if m.Cursor >= len(m.Types) {
				m.Cursor = 0
			}

		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Types) - 1
			}
		}
	}

	return m, nil
}

// View renders the current UI state as a string
func (m TypeSelectionModel) View() string {
	s := strings.Builder{}
	s.WriteString("Select a type:\n\n")

	for i, t := range m.Types {
		if m.Cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(string(t))
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

// NewTypeSelectionModel creates a new TypeSelectionModel
func NewTypeSelectionModel() TypeSelectionModel {
	return TypeSelectionModel{
		Types:  []BastionType{TypeBastion, TypeSession},
		Cursor: 0,
	}
}

// BastionModel represents the TUI model for bastion selection
// It implements the tea.Model interface required by Bubble Tea
type BastionModel struct {
	Cursor   int
	Choice   string
	Bastions []bastionSvc.Bastion
}

// Init initializes the model
// This is part of the tea.Model interface
func (m BastionModel) Init() tea.Cmd {
	// No initialization needed
	return nil
}

// Update handles messages and updates the model accordingly
// This is part of the tea.Model interface
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
			// Move the cursor down (with wrap-around)
			m.Cursor++
			if m.Cursor >= len(m.Bastions) {
				m.Cursor = 0
			}

		case "up", "k":
			// Move the cursor up (with wrap-around)
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Bastions) - 1
			}
		}
	}

	return m, nil
}

// View renders the current UI state as a string
// This is part of the tea.Model interface
func (m BastionModel) View() string {
	s := strings.Builder{}
	s.WriteString("Select a bastion host to create a session:\n\n")

	for i, bastion := range m.Bastions {
		if m.Cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(bastion.Name)
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

// NewBastionModel creates a new BastionModel with the provided bastions
func NewBastionModel(bastions []bastionSvc.Bastion) BastionModel {
	return BastionModel{
		Bastions: bastions,
		Cursor:   0,
	}
}

// SessionTypeModel represents the TUI model for session type selection
// It implements the tea.Model interface required by Bubble Tea
type SessionTypeModel struct {
	Cursor    int
	Choice    SessionType
	Types     []SessionType
	BastionID string
}

// Init initializes the model
// This is part of the tea.Model interface
func (m SessionTypeModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model accordingly
// This is part of the tea.Model interface
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
			// Move the cursor down (with wrap-around)
			m.Cursor++
			if m.Cursor >= len(m.Types) {
				m.Cursor = 0
			}

		case "up", "k":
			// Move the cursor up (with wrap-around)
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Types) - 1
			}
		}
	}

	return m, nil
}

// View renders the current UI state as a string
// This is part of the tea.Model interface
func (m SessionTypeModel) View() string {
	s := strings.Builder{}
	s.WriteString("Select a session type:\n\n")

	for i, t := range m.Types {
		if m.Cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(string(t))
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

// NewSessionTypeModel creates a new SessionTypeModel with the provided bastion ID
func NewSessionTypeModel(bastionID string) SessionTypeModel {
	return SessionTypeModel{
		Types:     []SessionType{TypeManagedSSH, TypePortForwarding},
		Cursor:    0,
		BastionID: bastionID,
	}
}

// TargetTypeModel represents the TUI model for target type selection
// It implements the tea.Model interface required by Bubble Tea
type TargetTypeModel struct {
	Cursor    int
	Choice    TargetType
	Types     []TargetType
	BastionID string
}

// Init initializes the model
// This is part of the tea.Model interface
func (m TargetTypeModel) Init() tea.Cmd {
	// No initialization needed
	return nil
}

// Update handles messages and updates the model accordingly
// This is part of the tea.Model interface
func (m TargetTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Set the choice and exit
			if m.Cursor >= 0 && m.Cursor < len(m.Types) {
				m.Choice = m.Types[m.Cursor]
			}
			return m, tea.Quit

		case "down", "j":
			// Move the cursor down (with wrap-around)
			m.Cursor++
			if m.Cursor >= len(m.Types) {
				m.Cursor = 0
			}

		case "up", "k":
			// Move the cursor up (with wrap-around)
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Types) - 1
			}
		}
	}

	return m, nil
}

// View renders the current UI state as a string
// This is part of the tea.Model interface
func (m TargetTypeModel) View() string {
	s := strings.Builder{}
	s.WriteString("Select a target type:\n\n")

	// Render each target type option with a cursor indicator
	for i, t := range m.Types {
		if m.Cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(string(t))
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

// NewTargetTypeModel creates a new TargetTypeModel with the provided bastion ID
func NewTargetTypeModel(bastionID string) TargetTypeModel {
	return TargetTypeModel{
		Types:     []TargetType{TargetOKE, TargetDatabase, TargetInstance},
		Cursor:    0,
		BastionID: bastionID,
	}
}

// OKEClusterModel allows selecting an OKE cluster
// Implements tea.Model
type OKEClusterModel struct {
	Cursor   int
	Choice   string
	Clusters []okeSvc.Cluster
}

func (m OKEClusterModel) Init() tea.Cmd { return nil }

func (m OKEClusterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			if m.Cursor >= 0 && m.Cursor < len(m.Clusters) {
				m.Choice = m.Clusters[m.Cursor].ID
			}
			return m, tea.Quit
		case "down", "j":
			m.Cursor++
			if m.Cursor >= len(m.Clusters) {
				m.Cursor = 0
			}
		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Clusters) - 1
			}
		}
	}
	return m, nil
}

func (m OKEClusterModel) View() string {
	var b strings.Builder
	b.WriteString("Select an OKE cluster:\n\n")
	for i, c := range m.Clusters {
		if m.Cursor == i {
			b.WriteString("(•) ")
		} else {
			b.WriteString("( ) ")
		}
		b.WriteString(c.Name)
		b.WriteString("\n")
	}
	b.WriteString("\n(press q to quit)\n")
	return b.String()
}

func NewOKEClusterModel(items []okeSvc.Cluster) OKEClusterModel {
	return OKEClusterModel{Clusters: items, Cursor: 0}
}

// InstanceSelectionModel allows selecting a Compute Instance
type InstanceSelectionModel struct {
	Cursor    int
	Choice    string
	Instances []instanceSvc.Instance
}

func (m InstanceSelectionModel) Init() tea.Cmd { return nil }
func (m InstanceSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			if m.Cursor >= 0 && m.Cursor < len(m.Instances) {
				m.Choice = m.Instances[m.Cursor].ID
			}
			return m, tea.Quit
		case "down", "j":
			m.Cursor++
			if m.Cursor >= len(m.Instances) {
				m.Cursor = 0
			}
		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Instances) - 1
			}
		}
	}
	return m, nil
}
func (m InstanceSelectionModel) View() string {
	var b strings.Builder
	b.WriteString("Select an Instance:\n\n")
	for i, inst := range m.Instances {
		if m.Cursor == i {
			b.WriteString("(•) ")
		} else {
			b.WriteString("( ) ")
		}
		name := inst.Name
		if name == "" {
			name = inst.Hostname
		}
		if name == "" {
			name = inst.ID
		}
		b.WriteString(name)
		b.WriteString("\n")
	}
	b.WriteString("\n(press q to quit)\n")
	return b.String()
}
func NewInstanceSelectionModel(items []instanceSvc.Instance) InstanceSelectionModel {
	return InstanceSelectionModel{Instances: items, Cursor: 0}
}

// DatabaseSelectionModel allows selecting an Autonomous Database
type DatabaseSelectionModel struct {
	Cursor    int
	Choice    string
	Databases []autonomousdbSvc.AutonomousDatabase
}

func (m DatabaseSelectionModel) Init() tea.Cmd { return nil }
func (m DatabaseSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			if m.Cursor >= 0 && m.Cursor < len(m.Databases) {
				m.Choice = m.Databases[m.Cursor].ID
			}
			return m, tea.Quit
		case "down", "j":
			m.Cursor++
			if m.Cursor >= len(m.Databases) {
				m.Cursor = 0
			}
		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Databases) - 1
			}
		}
	}
	return m, nil
}
func (m DatabaseSelectionModel) View() string {
	var b strings.Builder
	b.WriteString("Select an Autonomous Database:\n\n")
	for i, db := range m.Databases {
		if m.Cursor == i {
			b.WriteString("(•) ")
		} else {
			b.WriteString("( ) ")
		}
		name := db.Name
		if name == "" {
			name = db.ID
		}
		b.WriteString(name)
		b.WriteString("\n")
	}
	b.WriteString("\n(press q to quit)\n")
	return b.String()
}
func NewDatabaseSelectionModel(items []autonomousdbSvc.AutonomousDatabase) DatabaseSelectionModel {
	return DatabaseSelectionModel{Databases: items, Cursor: 0}
}

// Fancy searchable list implementation (for Instances/OKE/DBs)
// resourceItem implements list.Item

type resourceItem struct {
	id          string
	title       string
	description string
}

func (i resourceItem) Title() string       { return i.title }
func (i resourceItem) Description() string { return i.description }
func (i resourceItem) FilterValue() string { return i.title + " " + i.description }

// ResourceListModel wraps bubbles/list and stores selected id

type ResourceListModel struct {
	list   list.Model
	choice string
	keys   struct {
		confirm key.Binding
		quit    key.Binding
	}
}

func NewResourceListModel(title string, items []list.Item) ResourceListModel {
	delegate := list.NewDefaultDelegate()
	m := list.New(items, delegate, 0, 0)
	m.Title = title
	m.SetShowTitle(true)
	m.SetShowHelp(true)
	m.SetFilteringEnabled(true)
	m.SetShowFilter(true)

	rm := ResourceListModel{list: m}
	rm.keys.confirm = key.NewBinding(key.WithKeys("enter"))
	rm.keys.quit = key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"))
	return rm
}

// Helpers to build list models for each resource type
func NewInstanceListModelFancy(instances []instanceSvc.Instance) ResourceListModel {
	items := make([]list.Item, 0, len(instances))
	for _, inst := range instances {
		name := inst.Name
		if name == "" {
			name = inst.Hostname
		}
		if name == "" {
			name = inst.ID
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
	return NewResourceListModel("Instances", items)
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
	return NewResourceListModel("OKE Clusters", items)
}

func NewDBListModelFancy(dbs []autonomousdbSvc.AutonomousDatabase) ResourceListModel {
	items := make([]list.Item, 0, len(dbs))
	for _, d := range dbs {
		desc := d.PrivateEndpoint
		items = append(items, resourceItem{id: d.ID, title: d.Name, description: desc})
	}
	return NewResourceListModel("Autonomous Databases", items)
}

func (m ResourceListModel) Init() tea.Cmd { return nil }

func (m ResourceListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ResourceListModel) View() string { return m.list.View() }

// Choice returns selected ID
func (m ResourceListModel) Choice() string { return m.choice }
