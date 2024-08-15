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

// ApiResourceListState represents the current state of the ApiResource list view.
type ApiResourceListState int

const (
	StateNotStarted ApiResourceListState = iota
	StateFetching
	StateCompleted
	StateError
)

type ApiResourceListModel struct {
	styles        *tui.Styles
	spinner       spinner.Model
	width, height int
	cli           *core.CLI
	state         ApiResourceListState
	stateError    error
	list          list.Model
}

func NewApiResourceListModel(cli *core.CLI) ApiResourceListModel {

	return ApiResourceListModel{
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

func (m *ApiResourceListModel) fetchApiResources() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := m.cli.API.APIResource.List(ctx, "BUSINESS")
	if err != nil {
		return err
	}
	return list
}

// Init initializes the model and returns the initial command.
func (m ApiResourceListModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchApiResources,
		m.spinner.Tick,
	)
}

func (m ApiResourceListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case *models.APIResourceList:
		ApiResources := []list.Item{}
		for _, app := range msg.APIResources {
			ApiResources = append(ApiResources, tui.NewItem(app.Name, app.ID))
		}
		m.list = list.New(ApiResources, list.NewDefaultDelegate(), 0, 0)
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

func (m ApiResourceListModel) View() string {
	switch m.state {
	case StateFetching:
		return fmt.Sprintf("\n\n   %s Fetching Api Resources...!\n\n", m.spinner.View())
	case StateCompleted:
		m.list.Title = "Api Resources"
		return m.styles.List.Render(m.list.View())
	case StateError:
		return fmt.Sprint(m.stateError.Error())
	}
	return ""
}
