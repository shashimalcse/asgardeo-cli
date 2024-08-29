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
	StateFetching ApiResourceListState = iota
	StateCompleted
	StateError
	StateApiResourceSelected
)

type ApiResourceListModel struct {
	styles              *tui.Styles
	spinner             spinner.Model
	width, height       int
	cli                 *core.CLI
	state               ApiResourceListState
	stateError          error
	list                list.Model
	scopeList           list.Model
	apiResourceList     []models.APIResource
	selectedApiResource models.APIResource
}

func NewApiResourceListModel(cli *core.CLI) *ApiResourceListModel {

	return &ApiResourceListModel{
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

	apiList, err := m.cli.API.APIResource.List(ctx, "BUSINESS")
	if err != nil {
		return err
	}
	return apiList
}

func (m *ApiResourceListModel) fetchApiResource() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	apiResource, err := m.cli.API.APIResource.Get(ctx, m.selectedApiResource.ID)
	if err != nil {
		return err
	}
	m.selectedApiResource = *apiResource

	return nil
}

// Init initializes the model and returns the initial command.
func (m *ApiResourceListModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchApiResources,
		m.spinner.Tick,
	)
}

func (m *ApiResourceListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			selectedItem := m.list.SelectedItem()
			if selectedItem == nil {
				return m, nil
			}
			m.selectedApiResource = m.getApiResourceByName(selectedItem.FilterValue())
			err := m.fetchApiResource()
			if err != nil {
				m.state = StateError
				m.stateError = err
				return m, tea.Quit
			}
			var scopes []list.Item
			for _, scope := range m.selectedApiResource.Scopes {
				scopes = append(scopes, tui.NewItem(scope.Name, scope.ID))
			}
			m.scopeList = list.New(scopes, list.NewDefaultDelegate(), 0, 0)
			m.scopeList.Title = "API Resource : " + m.selectedApiResource.Name + " Scopes"
			h, v := m.styles.List.GetFrameSize()
			m.scopeList.SetSize(m.width-h, m.height-v)
			m.state = StateApiResourceSelected
		}
	case *models.APIResourceList:
		m.apiResourceList = msg.APIResources
		var ApiResources []list.Item
		for _, app := range msg.APIResources {
			ApiResources = append(ApiResources, tui.NewItem(app.Name, app.ID))
		}
		m.list = list.New(ApiResources, list.NewDefaultDelegate(), 0, 0)
		m.list.Title = "Api Resources"
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
	if m.state == StateApiResourceSelected {
		m.scopeList, _ = m.scopeList.Update(msg)
	}
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m *ApiResourceListModel) View() string {
	switch m.state {
	case StateFetching:
		return fmt.Sprintf("\n\n   %s Fetching Api Resources...!\n\n", m.spinner.View())
	case StateCompleted:
		m.list.Title = "Api Resources"
		return m.styles.List.Render(m.list.View())
	case StateApiResourceSelected:
		return m.styles.List.Render(m.scopeList.View())
	case StateError:
		return fmt.Sprint(m.stateError.Error())
	}
	return ""
}

func (m *ApiResourceListModel) getApiResourceByName(name string) models.APIResource {
	for _, api := range m.apiResourceList {
		if api.Name == name {
			return api
		}
	}
	return models.APIResource{}
}
