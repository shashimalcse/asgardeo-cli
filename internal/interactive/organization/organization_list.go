package interactive

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shashimalcse/asgardeo-cli/internal/core"
	"github.com/shashimalcse/asgardeo-cli/internal/models"
	"github.com/shashimalcse/asgardeo-cli/internal/tui"
)

// OrganizationListState represents the current state of the organization list view.
type OrganizationListState int

const (
	StateFetching OrganizationListState = iota
	StateCompleted
	StateError
)

type OrganizationListModel struct {
	styles        *tui.Styles
	spinner       spinner.Model
	width, height int
	cli           *core.CLI
	state         OrganizationListState
	stateError    error
	list          list.Model
	filter        string
}

func NewOrganizationListModel(cli *core.CLI, filter string) *OrganizationListModel {
	return &OrganizationListModel{
		styles:  tui.DefaultStyles(),
		spinner: newSpinner(),
		cli:     cli,
		state:   StateFetching,
		filter:  filter,
	}
}

func newSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

func (m *OrganizationListModel) fetchOrganizations() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	orgsList, err := m.cli.API.Organization.List(ctx, m.filter)
	if err != nil {
		return err
	}
	return orgsList
}

// Init initializes the model and returns the initial command.
func (m *OrganizationListModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchOrganizations,
		m.spinner.Tick,
	)
}

func (m *OrganizationListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || (msg.Type == tea.KeyRunes && msg.String() == "q") {
			return m, tea.Quit
		}
	case *models.OrganizationList:
		var organizations []list.Item
		for _, org := range msg.Organizations {
			organizations = append(organizations, tui.NewItem(org.Name, org.ID))
		}
		m.list = list.New(organizations, list.NewDefaultDelegate(), 0, 0)
		h, v := m.styles.List.GetFrameSize()
		m.list.SetSize(m.width-h, m.height-v)
		m.state = StateCompleted
		return m, nil
	case error:
		m.state = StateError
		m.stateError = msg
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	if m.state == StateCompleted {
		m.list, _ = m.list.Update(msg)
	}
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m *OrganizationListModel) View() string {
	switch m.state {
	case StateFetching:
		return fmt.Sprintf("\n\n   %s Fetching organizations...!\n\n", m.spinner.View())
	case StateCompleted:
		m.list.Title = "Organizations"
		return m.styles.List.Render(m.list.View())
	case StateError:
		return fmt.Sprint(m.stateError.Error())
	}
	return ""
}
