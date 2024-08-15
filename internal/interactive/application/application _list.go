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

// ApplicationListState represents the current state of the application list view.
type ApplicationListState int

const (
	StateNotStarted ApplicationListState = iota
	StateFetching
	StateCompleted
	StateError
)

type ApplicationListModel struct {
	styles        *tui.Styles
	spinner       spinner.Model
	width, height int
	cli           *core.CLI
	state         ApplicationListState
	stateError    error
	list          list.Model
}

func NewApplicationListModel(cli *core.CLI) ApplicationListModel {

	return ApplicationListModel{
		styles:  tui.DefaultStyles(),
		spinner: newSpinner(),
		cli:     cli,
		state:   StateFetching,
	}

}

func newSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

func (m *ApplicationListModel) fetchApplications() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := m.cli.API.Application.List(ctx)
	if err != nil {
		return err
	}
	return list
}

// Init initializes the model and returns the initial command.
func (m ApplicationListModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchApplications,
		m.spinner.Tick,
	)
}

func (m ApplicationListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case *models.ApplicationList:
		applications := []list.Item{}
		for _, app := range msg.Applications {
			applications = append(applications, tui.NewItem(app.Name, app.ID))
		}
		m.list = list.New(applications, list.NewDefaultDelegate(), 0, 0)
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

func (m ApplicationListModel) View() string {
	switch m.state {
	case StateFetching:
		return fmt.Sprintf("\n\n   %s Fetching applications...!\n\n", m.spinner.View())
	case StateCompleted:
		m.list.Title = "Applications"
		return m.styles.List.Render(m.list.View())
	case StateError:
		return fmt.Sprint(m.stateError.Error())
	}
	return ""
}
