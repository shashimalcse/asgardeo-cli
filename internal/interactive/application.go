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

// Application List Interactive
type ApplicationListState int

const (
	ApplicationListFetchingNotStarted ApplicationListState = iota
	ApplicationListFetchingInProgress
	ApplicationListFetchingCompleted
	ApplicationListFetchingError
)

type ApplicationListModel struct {
	styles     *tui.Styles
	spinner    spinner.Model
	width      int
	height     int
	cli        *core.CLI
	state      ApplicationListState
	stateError error
	list       list.Model
}

func NewApplicationListModel(cli *core.CLI) ApplicationListModel {

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return ApplicationListModel{
		styles:  tui.DefaultStyles(),
		spinner: s,
		cli:     cli,
		state:   ApplicationListFetchingInProgress,
	}

}

func (m ApplicationListModel) fetchApplications() tea.Cmd {
	return func() tea.Msg {
		list, err := m.cli.API.Application.List(context.Background())
		if err != nil {
			m.state = ApplicationListFetchingError
			return err
		}
		return list
	}
}

func (m ApplicationListModel) Init() tea.Cmd {
	return tea.Batch(m.fetchApplications(), m.spinner.Tick)
}

func (m ApplicationListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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

func (m ApplicationListModel) View() string {
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

// Application Create Interactive

type ApplicationCreateState int

const (
	ApplicationCreateInitiated ApplicationCreateState = iota
	ApplicationCreateTypeSelected
	ApplicationCreateQuestionsCompleted
	ApplicationCreateFetchingError
)

type ApplicationType string

const (
	SinglePage  ApplicationType = "Single-Page Application"
	Traditional ApplicationType = "Traditional Web Application"
	Mobile      ApplicationType = "Mobile Application"
	Standard    ApplicationType = "Standard-Based Application"
	M2M         ApplicationType = "M2M Application"
)

type ApplicationCreateModel struct {
	styles                         *tui.Styles
	spinner                        spinner.Model
	width                          int
	height                         int
	cli                            *core.CLI
	state                          ApplicationCreateState
	stateError                     error
	applicationTypesList           list.Model
	questionsForSinglePage         []tui.Question
	currentSinglePageQuestionIndex int
	applicationType                ApplicationType
}

func NewApplicationCreateModel(cli *core.CLI) ApplicationCreateModel {

	applicationTypesItems := []list.Item{
		tui.NewItemWithKey("single_page", "Single-Page Application", "A web application that runs application logic in the browser."),
		tui.NewItemWithKey("traditional", "Traditional Web Application", "A web application that runs application logic on the server."),
		tui.NewItemWithKey("mobile", "Mobile Application", "Applications developed to target mobile devices."),
		tui.NewItemWithKey("standard", "Standard-Based Application", "Applications built using standard protocols."),
		tui.NewItemWithKey("m2m", "M2M Application", "Applications tailored for Machine to Machine communication."),
	}
	applicationTypesList := list.New(applicationTypesItems, list.NewDefaultDelegate(), 0, 0)
	applicationTypesList.Title = "Select application template to create application"

	questionsForSinglePage := []tui.Question{
		tui.NewQuestion("Name", "Name", tui.ShortQuestion),
		tui.NewQuestion("Authorized redirect URL", "Authorized redirect URL", tui.ShortQuestion),
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return ApplicationCreateModel{
		styles:                         tui.DefaultStyles(),
		spinner:                        s,
		cli:                            cli,
		state:                          ApplicationCreateInitiated,
		applicationTypesList:           applicationTypesList,
		questionsForSinglePage:         questionsForSinglePage,
		currentSinglePageQuestionIndex: 0,
	}

}

func (m ApplicationCreateModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m ApplicationCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	currentSinglePageQuestion := m.questionsForSinglePage[m.currentSinglePageQuestionIndex]
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.state {
			case ApplicationCreateInitiated:
				i, ok := m.applicationTypesList.SelectedItem().(tui.Item)
				if ok {
					m.applicationType = ApplicationType(i.Title())
					m.state = ApplicationCreateTypeSelected
				}
			case ApplicationCreateTypeSelected:
				switch m.applicationType {
				case SinglePage:
					currentSinglePageQuestion := &m.questionsForSinglePage[m.currentSinglePageQuestionIndex]
					currentSinglePageQuestion.Answer = currentSinglePageQuestion.Input.Value()
					if m.currentSinglePageQuestionIndex == len(m.questionsForSinglePage)-1 {
						m.state = ApplicationCreateQuestionsCompleted
					} else {
						m.NextSinglePageQuestion()
					}
					return m, currentSinglePageQuestion.Input.Blur
				}
			}

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := m.styles.List.GetFrameSize()
		m.applicationTypesList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.applicationTypesList, _ = m.applicationTypesList.Update(msg)
	currentSinglePageQuestion.Input, _ = currentSinglePageQuestion.Input.Update(msg)
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m ApplicationCreateModel) View() string {

	switch m.state {
	case ApplicationCreateInitiated:
		return m.styles.List.Render(m.applicationTypesList.View())
	case ApplicationCreateTypeSelected:
		switch m.applicationType {
		case SinglePage:
			current := m.questionsForSinglePage[m.currentSinglePageQuestionIndex]
			var previousQAs string
			for i := 0; i < m.currentSinglePageQuestionIndex; i++ {
				question := m.questionsForSinglePage[i]
				previousQAs += fmt.Sprintf("%s : %s\n", question.Question, question.Answer)
			}
			return previousQAs + current.Input.View()
		default:
			return "Other types are not supported yet!"
		}

	case ApplicationCreateQuestionsCompleted:
		var previousQAs string
		for i := 0; i < len(m.questionsForSinglePage); i++ {
			question := m.questionsForSinglePage[i]
			previousQAs += fmt.Sprintf("%s : %s\n", question.Question, question.Answer)
		}
		return previousQAs
	}

	return ""
}

func (m *ApplicationCreateModel) NextSinglePageQuestion() {
	if m.currentSinglePageQuestionIndex < len(m.questionsForSinglePage)-1 {
		m.currentSinglePageQuestionIndex++
	} else {
		m.currentSinglePageQuestionIndex = 0
	}
}
