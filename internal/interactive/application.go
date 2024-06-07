package interactive

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shashimalcse/is-cli/internal/core"
	"github.com/shashimalcse/is-cli/internal/management"
	"github.com/shashimalcse/is-cli/internal/tui"
)

// AuthenticateState represents the state of the authentication process for both machine and user
type ApplicationListState int

const (
	ApplicationListFetchingNotStarted ApplicationListState = iota
	ApplicationListFetchingInProgress
	ApplicationListFetchingCompleted
	ApplicationListFetchingError
)

type ApplicationModel struct {
	styles     *tui.Styles
	spinner    spinner.Model
	width      int
	height     int
	cli        *core.CLI
	state      ApplicationListState
	stateError error
	list       list.Model
}

func NewApplicationModel(cli *core.CLI) ApplicationModel {

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return ApplicationModel{
		styles:  tui.DefaultStyles(),
		spinner: s,
		cli:     cli,
		state:   ApplicationListFetchingInProgress,
	}

}

func (m ApplicationModel) fetchApplications() tea.Cmd {
	return func() tea.Msg {
		list, err := m.cli.API.Application.List(context.Background())
		if err != nil {
			m.state = ApplicationListFetchingError
			return err
		}
		return list
	}
}

func (m ApplicationModel) Init() tea.Cmd {
	return tea.Batch(m.fetchApplications(), m.spinner.Tick)
}

func (m ApplicationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case *management.ApplicationList:
		applications := []list.Item{}
		for _, app := range msg.Applications {
			applications = append(applications, tui.NewItem(app.Name, app.ID))
		}
		m.list = list.New(applications, list.NewDefaultDelegate(), 0, 0)
		h, v := m.styles.List.GetFrameSize()
		m.list.SetSize(m.width-h, m.height-v)
		m.state = ApplicationListFetchingCompleted
		return m, nil
	case error:
		m.state = ApplicationListFetchingError
		m.stateError = msg
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	if m.state == ApplicationListFetchingCompleted {
		m.list, _ = m.list.Update(msg)
	}
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m ApplicationModel) View() string {
	switch m.state {
	case ApplicationListFetchingInProgress:
		return fmt.Sprintf("\n\n   %s Fetching applications...!\n\n", m.spinner.View())
	case ApplicationListFetchingCompleted:
		return m.styles.List.Render(m.list.View())
	case ApplicationListFetchingError:
		return fmt.Sprint(m.stateError.Error())
	}
	return ""
}
