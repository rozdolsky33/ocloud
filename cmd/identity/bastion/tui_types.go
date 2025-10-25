package bastion

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	instSvc "github.com/rozdolsky33/ocloud/internal/services/compute/instance"
	okeSvc "github.com/rozdolsky33/ocloud/internal/services/compute/oke"
	adbSvc "github.com/rozdolsky33/ocloud/internal/services/database/autonomousdb"
	hwdbSvc "github.com/rozdolsky33/ocloud/internal/services/database/heatwavedb"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
)

// BastionType identifies the top-level action.
type BastionType string

const (
	TypeBastion BastionType = "Bastion"
	TypeSession BastionType = "Session"
)

// TargetType identifies what the session connects to.
type TargetType string

const (
	TargetOKE      TargetType = "OKE"
	TargetDatabase TargetType = "Database"
	TargetInstance TargetType = "Instance"
)

// SessionType identifies how the bastion session behaves.
type SessionType string

const (
	TypeManagedSSH     SessionType = "Managed SSH"
	TypePortForwarding SessionType = "Port-Forwarding"
)

// DatabaseType identifies the type of database to connect to.
type DatabaseType string

const (
	DatabaseHeatWave   DatabaseType = "MySQL HeatWave"
	DatabaseAutonomous DatabaseType = "Autonomous Database"
)

//-----------------------------------Bastion/Session Creation Selection-------------------------------------------------

// TypeSelectionModel defines a TUI model for selecting a BastionType from a list of available types.
type TypeSelectionModel struct {
	Cursor int
	Choice BastionType
	Types  []BastionType
}

// NewTypeSelectionModel creates a new TypeSelectionModel.
func NewTypeSelectionModel() TypeSelectionModel {
	return TypeSelectionModel{
		Types:  []BastionType{TypeBastion, TypeSession},
		Cursor: 0,
	}
}

// Init initializes the TypeSelectionModel and returns a command.
func (m TypeSelectionModel) Init() tea.Cmd { return nil }

// Update processes input messages to update the model's state and returns the updated model and an optional command.
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

// View returns a string representation of the model's state.
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

//--------------------------------------------------Bastion Selection---------------------------------------------------

// BastionModel Bastion Selection
type BastionModel struct {
	Cursor   int
	Choice   string
	Bastions []bastionSvc.Bastion
}

// NewBastionModel creates a BastionModel instance with the provided list of bastions and initializes the cursor to 0.
func NewBastionModel(bastions []bastionSvc.Bastion) BastionModel {
	return BastionModel{Bastions: bastions, Cursor: 0}
}

// Init initializes the BastionModel and returns an optional command to execute.
func (m BastionModel) Init() tea.Cmd { return nil }

// Update processes incoming messages, updates the model's state, and determines the next command to execute.
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

// View renders the string representation of the BastionModel, displaying the list of bastion hosts and current selection.
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

//------------------------------------------Session Type Selection------------------------------------------------------

// SessionTypeModel Session Type Selection
type SessionTypeModel struct {
	Cursor    int
	Choice    SessionType
	Types     []SessionType
	BastionID string
}

// NewSessionTypeModel creates a SessionTypeModel instance with the provided list of bastions and initializes the cursor to 0.
func NewSessionTypeModel(bastionID string) SessionTypeModel {
	return SessionTypeModel{
		Types:     []SessionType{TypeManagedSSH, TypePortForwarding},
		Cursor:    0,
		BastionID: bastionID,
	}
}

// Init initializes the SessionTypeModel and returns an optional command to execute.
func (m SessionTypeModel) Init() tea.Cmd { return nil }

// Update processes incoming messages, updates the model's state, and determines the next command to execute.
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

// View returns a string representation of the model's state.
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

//------------------------------------------Database Type Selection-----------------------------------------------------

// DatabaseTypeModel Database Type Selection
type DatabaseTypeModel struct {
	Cursor int
	Choice DatabaseType
	Types  []DatabaseType
}

// NewDatabaseTypeModel creates a DatabaseTypeModel instance with HeatWave and Autonomous Database options.
func NewDatabaseTypeModel() DatabaseTypeModel {
	return DatabaseTypeModel{
		Types:  []DatabaseType{DatabaseHeatWave, DatabaseAutonomous},
		Cursor: 0,
	}
}

// Init initializes the DatabaseTypeModel and returns an optional command to execute.
func (m DatabaseTypeModel) Init() tea.Cmd { return nil }

// Update processes incoming messages, updates the model's state, and determines the next command to execute.
func (m DatabaseTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

// View returns a string representation of the model's state.
func (m DatabaseTypeModel) View() string {
	var b strings.Builder
	b.WriteString("Select database type:\n\n")
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

//------------------------------------------Target Type Selection-------------------------------------------------------

// TargetTypeModel Target Type Selection
type TargetTypeModel struct {
	Cursor    int
	Choice    TargetType
	Types     []TargetType
	BastionID string
}

// NewTargetTypeModel creates a TargetTypeModel instance with the provided list of bastions and initializes the cursor to 0.
func NewTargetTypeModel(bastionID string) TargetTypeModel {
	var types []TargetType
	types = []TargetType{TargetOKE, TargetDatabase, TargetInstance}
	return TargetTypeModel{Types: types, Cursor: 0, BastionID: bastionID}
}

// Init initializes the TargetTypeModel and returns an optional command to execute.
func (m TargetTypeModel) Init() tea.Cmd { return nil }

// Update processes incoming messages, updates the model's state, and determines the next command to execute.
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

// View returns a string representation of the model's state.
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

// ----------------------------------Fancy searchable list (Instances / OKE / DB)--------------------------------------
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

// NewInstanceListModelFancy creates a ResourceListModel to display instances in a searchable and interactive list.
func NewInstanceListModelFancy(instances []instSvc.Instance) ResourceListModel {
	items := make([]list.Item, 0, len(instances))
	for _, inst := range instances {
		name := inst.DisplayName
		if name == "" {
			name = inst.OCID
		}
		desc := fmt.Sprintf("IP: %s", inst.PrimaryIP)
		if inst.VcnName != "" {
			desc = inst.VcnName
			if inst.SubnetName != "" {
				desc += " · " + inst.SubnetName
			}
		}
		items = append(items, resourceItem{id: inst.OCID, title: name, description: desc})
	}
	return newResourceList("Instances", items)
}

// NewOKEListModelFancy creates a ResourceListModel to display OKE clusters in a searchable and interactive list.
func NewOKEListModelFancy(clusters []okeSvc.Cluster) ResourceListModel {
	items := make([]list.Item, 0, len(clusters))
	for _, c := range clusters {
		desc := c.KubernetesVersion
		if c.PrivateEndpoint != "" {
			desc += " · PE"
		}
		items = append(items, resourceItem{id: c.OCID, title: c.DisplayName, description: desc})
	}
	return newResourceList("OKE Clusters", items)
}

// NewDBListModelFancy creates a ResourceListModel populated with a list of autonomous databases for TUI display.
func NewDBListModelFancy(dbs []adbSvc.AutonomousDatabase) ResourceListModel {
	items := make([]list.Item, 0, len(dbs))
	for _, d := range dbs {
		desc := d.PrivateEndpoint
		items = append(items, resourceItem{id: d.ID, title: d.Name, description: desc})
	}
	return newResourceList("Autonomous Databases", items)
}

// NewHeatWaveDBListModelFancy creates a ResourceListModel populated with a list of HeatWave databases for TUI display.
func NewHeatWaveDBListModelFancy(dbs []hwdbSvc.HeatWaveDatabase) ResourceListModel {
	items := make([]list.Item, 0, len(dbs))
	for _, d := range dbs {
		desc := d.IpAddress
		items = append(items, resourceItem{id: d.ID, title: d.DisplayName, description: desc})
	}
	return newResourceList("HeatWave Databases", items)
}

//---------------------------------------SSH Keys----------------------------------------------------------------------

// SSHFileItem is a list item representing a file system entry (file or directory).
type SSHFileItem struct {
	path       string
	title      string
	permission string
	isDir      bool
}

// Title returns the display title (implements list.Item).
func (i SSHFileItem) Title() string { return i.title }

// Description returns permissions or metadata (implements list.Item).
func (i SSHFileItem) Description() string { return i.permission }
func (i SSHFileItem) FilterValue() string { return i.title + " " + i.permission }

// SSHFilesModel is the canonical model name for SSH file selection/browsing.
type SSHFilesModel struct {
	list       list.Model
	choice     string
	currentDir string
	showPublic bool
	browsing   bool
	keys       struct {
		confirm key.Binding
		quit    key.Binding
		upDir   key.Binding
	}
}

// SHHFilesModel is an alias for SSHFilesModel used in contexts requiring SSH file selection and interaction.
type SHHFilesModel = SSHFilesModel

func (m SSHFilesModel) Init() tea.Cmd { return nil }
func (m SSHFilesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.quit) {
			return m, tea.Quit
		}
		if m.browsing && key.Matches(msg, m.keys.upDir) {
			if m.currentDir != "" {
				parent := path.Dir(m.currentDir)
				if parent != m.currentDir {
					m.currentDir = parent
					m.NewSSHFilesModelFancyList()
				}
			}
			return m, nil
		}
		if key.Matches(msg, m.keys.confirm) {
			if it, ok := m.list.SelectedItem().(SSHFileItem); ok {
				if m.browsing && it.isDir {
					if it.path == ".." {
						parent := path.Dir(m.currentDir)
						if parent != m.currentDir {
							m.currentDir = parent
							m.NewSSHFilesModelFancyList()
						}
					} else {
						m.currentDir = it.path
						m.NewSSHFilesModelFancyList()
					}
					return m, nil
				}
				m.choice = it.path
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SSHFilesModel) View() string   { return m.list.View() }
func (m SSHFilesModel) Choice() string { return m.choice }

func newSSHList(title string, items []list.Item) SSHFilesModel {
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = title
	l.SetShowTitle(true)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)

	rm := SSHFilesModel{list: l}
	rm.keys.confirm = key.NewBinding(key.WithKeys("enter"))
	rm.keys.quit = key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"))
	rm.keys.upDir = key.NewBinding(key.WithKeys("backspace", "left"))
	return rm
}

// filePermString returns the file's unix permission bits as a short octal string (e.g., "600").
func filePermString(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "n/a"
	}
	perm := info.Mode().Perm()
	return fmt.Sprintf("%o", perm)
}

// NewSSHFilesModelFancyList populates the list items based on currentDir and filtering rules.
func (m *SSHFilesModel) NewSSHFilesModelFancyList() {
	if !m.browsing || m.currentDir == "" {
		return
	}
	entries, err := os.ReadDir(m.currentDir)
	if err != nil {
		return
	}
	items := make([]list.Item, 0, len(entries)+1)
	if parent := path.Dir(m.currentDir); parent != m.currentDir {
		items = append(items, SSHFileItem{path: "..", title: "..", permission: "", isDir: true})
	}
	for _, e := range entries {
		if e.IsDir() {
			p := path.Join(m.currentDir, e.Name())
			items = append(items, SSHFileItem{path: p, title: e.Name() + string(os.PathSeparator), permission: "dir", isDir: true})
		}
	}
	for _, e := range entries {
		if !e.IsDir() {
			name := e.Name()
			if m.showPublic && !strings.HasSuffix(name, ".pub") {
				continue
			}
			if !m.showPublic && strings.HasSuffix(name, ".pub") {
				continue
			}
			p := path.Join(m.currentDir, name)
			items = append(items, SSHFileItem{path: p, title: name, permission: filePermString(p), isDir: false})
		}
	}
	m.list.SetItems(items)
}

// NewSSHKeysModelBrowser creates a navigable SSHFilesModel starting from startDir.
func NewSSHKeysModelBrowser(title, startDir string, showPublic bool) SHHFilesModel {
	m := newSSHList(title, nil)
	m.browsing = true
	m.currentDir = startDir
	m.showPublic = showPublic
	m.NewSSHFilesModelFancyList()
	return m
}
