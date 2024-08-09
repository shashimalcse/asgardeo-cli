package interactive

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
	"github.com/shashimalcse/is-cli/internal/auth"
	"github.com/shashimalcse/is-cli/internal/core"
	"github.com/shashimalcse/is-cli/internal/models"
	"github.com/shashimalcse/is-cli/internal/tui"
)

var (
	// Login options
	AS_A_MACHINE = "as a machine"
	AS_A_USER    = "as a user"
)

// AuthenticateState represents the current state of the authentication process
type AuthenticateState int

// Constants for different authentication states
const (
	StateNotStarted AuthenticateState = iota
	StateClientCredentialsInProgress
	StateClientCredentialsCompleted
	StateClientCredentialsError
	StateDeviceFlowInitiated
	StateDeviceFlowCodeReceived
	StateDeviceFlowError
	StateDeviceFlowBrowserWait
	StateDeviceFlowBrowserCompleted
	StateDeviceFlowBrowserError
	StateDeviceFlowCompleted
)

type LoginModel struct {
	styles              *tui.Styles
	spinner             spinner.Model
	width, height       int
	loginOptions        list.Model
	isLoginOptionChosen bool
	loginOptionChosen   string
	questions           []tui.Question
	currentQuestionIdx  int
	questionsDone       bool
	cli                 *core.CLI
	state               AuthenticateState
	stateMessage        string
	deviceFlowState     auth.State
	outputResult        models.OutputResult
}

// NewLoginModel creates and initializes a new LoginModel
func NewLoginModel(cli *core.CLI) LoginModel {
	return LoginModel{
		styles:       tui.DefaultStyles(),
		spinner:      newSpinner(),
		loginOptions: newLoginOptions(),
		cli:          cli,
		state:        StateNotStarted,
	}
}

func newSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

func newLoginOptions() list.Model {
	items := []list.Item{
		tui.NewItem(AS_A_MACHINE, "Authenticates the IS CLI as a machine using client credentials"),
		tui.NewItem(AS_A_USER, "Authenticates the IS CLI as a user using personal credentials"),
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "How would you like to authenticate?"
	return l
}

func (m LoginModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m.handleKeyEnter(msg)
		}
	case tea.WindowSizeMsg:
		return m.handleWindowResize(msg)
	}

	var cmd tea.Cmd
	if m.isLoginOptionChosen && !m.questionsDone {
		m.questions[m.currentQuestionIdx].Input, _ = m.questions[m.currentQuestionIdx].Input.Update(msg)
	} else if !m.isLoginOptionChosen {
		m.loginOptions, _ = m.loginOptions.Update(msg)
	}
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m LoginModel) handleKeyEnter(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !m.isLoginOptionChosen {
		i, ok := m.loginOptions.SelectedItem().(tui.Item)
		if ok {
			m.loginOptionChosen = i.Title()
			m.isLoginOptionChosen = true
			m.initQuestions()
		}
	} else {
		currentQuestion := &m.questions[m.currentQuestionIdx]
		if m.loginOptionChosen == AS_A_USER {
			// Handle device flow
			if m.state == StateDeviceFlowBrowserWait {
				m.state = StateDeviceFlowBrowserCompleted
				err := m.getAccessTokenFromDeviceCode(m.deviceFlowState)
				if err != nil {
					m.state = StateDeviceFlowError
					m.stateMessage = err.Error()
				} else {
					m.state = StateDeviceFlowCompleted
					return m, tea.Quit
				}
			} else {
				if m.currentQuestionIdx == len(m.questions)-1 {
					m.questionsDone = true
					currentQuestion.Answer = currentQuestion.Input.Value()
					m.state = StateDeviceFlowInitiated
					state, err := m.getDeviceCode()
					if err != nil {
						m.state = StateDeviceFlowError
						m.stateMessage = err.Error()
						m.outputResult = models.OutputResult{
							Message: "Error getting device code: " + err.Error(),
							IsError: true,
						}
						return m, tea.Quit
					} else {
						m.state = StateDeviceFlowCodeReceived
						if err = browser.OpenURL(state.VerificationURIComplete); err != nil {
							m.state = StateDeviceFlowBrowserError
						}
						m.deviceFlowState = state
						m.state = StateDeviceFlowBrowserWait
					}
				} else {
					m.NextQuestion()
				}
				currentQuestion.Answer = currentQuestion.Input.Value()
				return m, currentQuestion.Input.Blur
			}
		} else {
			// Handle client credentials flow
			if m.currentQuestionIdx == len(m.questions)-1 {
				m.questionsDone = true
				currentQuestion.Answer = currentQuestion.Input.Value()
				m.state = StateClientCredentialsInProgress
				err := m.runLoginAsMachine()
				if err != nil {
					m.state = StateClientCredentialsError
					m.stateMessage = err.Error()
					m.outputResult = models.OutputResult{
						Message: err.Error(),
						IsError: true,
					}
				} else {
					m.state = StateClientCredentialsCompleted
					m.outputResult = models.OutputResult{
						Message: "Successfully authenticated as a machine",
						IsError: false,
					}
				}
				return m, tea.Quit
			} else {
				m.NextQuestion()
			}
			currentQuestion.Answer = currentQuestion.Input.Value()
			return m, currentQuestion.Input.Blur
		}
		currentQuestion.Input, _ = currentQuestion.Input.Update(msg)
	}
	return m, nil
}

func (m LoginModel) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width, m.height = msg.Width, msg.Height
	h, v := m.styles.List.GetFrameSize()
	m.loginOptions.SetSize(msg.Width-h, msg.Height-v)
	return m, nil
}

func (m *LoginModel) initQuestions() {
	if m.loginOptionChosen == AS_A_MACHINE {
		m.questions = []tui.Question{
			tui.NewQuestion("tenant", "Tenant Domain", tui.ShortQuestion),
			tui.NewQuestion("client id", "Client ID", tui.ShortQuestion),
			tui.NewQuestion("client secret", "Client Secret", tui.ShortSecretQuestion),
		}
	} else {
		m.questions = []tui.Question{
			tui.NewQuestion("tenant", "Tenant Domain", tui.ShortQuestion),
			tui.NewQuestion("client id", "Client ID", tui.ShortQuestion),
		}
	}
}

// View renders the current view of the model
func (m LoginModel) View() string {
	if !m.isLoginOptionChosen {
		return m.styles.List.Render(m.loginOptions.View())
	}

	if !m.questionsDone {
		var previousQAs string
		previousQAs += fmt.Sprintf("Trying to authenticate %s\n\n", m.loginOptionChosen)
		for i := 0; i < m.currentQuestionIdx; i++ {
			question := m.questions[i]
			previousQAs += fmt.Sprintf("%s : %s\n", question.Question, question.Answer)
		}
		return previousQAs + m.questions[m.currentQuestionIdx].Input.View()
	}

	return m.renderAuthenticationStatus()
}

func (m LoginModel) renderAuthenticationStatus() string {
	switch m.state {
	case StateClientCredentialsCompleted:
		return "Successfully authenticated as a machine"
	case StateClientCredentialsError:
		return "Error authenticating as a machine: " + m.stateMessage
	case StateClientCredentialsInProgress:
		return fmt.Sprintf("%s Authenticating as a machine...", m.spinner.View())
	case StateDeviceFlowInitiated:
		return "Device flow initiated."
	case StateDeviceFlowBrowserWait:
		return fmt.Sprintf("\n\n   %s Waiting for the login to complete in the browser. Press Enter after login is completed.\n\n", m.spinner.View())
	case StateDeviceFlowBrowserCompleted:
		return "Device flow completed."
	case StateDeviceFlowBrowserError:
		return "Error opening browser. Please visit " + m.deviceFlowState.VerificationURIComplete + " to authenticate."
	case StateDeviceFlowCompleted:
		return "Successfully logged in"
	case StateDeviceFlowError:
		return "Error initiating device flow: " + m.stateMessage
	default:
		return ""
	}
}

func (m *LoginModel) NextQuestion() {
	if m.currentQuestionIdx < len(m.questions)-1 {
		m.currentQuestionIdx++
	} else {
		m.currentQuestionIdx = 0
	}
}

func (m LoginModel) runLoginAsMachine() error {

	err := core.RunLoginAsMachine(
		core.LoginInputs{
			Tenant:       m.questions[0].Answer,
			ClientID:     m.questions[1].Answer,
			ClientSecret: m.questions[2].Answer,
		}, m.cli)
	return err
}

func (m LoginModel) getDeviceCode() (auth.State, error) {

	return core.GetDeviceCode(m.cli)
}

func (m LoginModel) getAccessTokenFromDeviceCode(state auth.State) error {

	return core.GetAccessTokenFromDeviceCode(m.cli, state)
}

func (m LoginModel) GetOutputValue() models.OutputResult {
	return m.outputResult
}
