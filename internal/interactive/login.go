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
	"github.com/shashimalcse/is-cli/internal/tui"
)

// AuthenticateState represents the state of the authentication process for both machine and user
type AuthenticateState int

const (
	NotStarted                  AuthenticateState = 0
	ClientCredentialsInProgress AuthenticateState = 1
	ClientCredentialsCompleted  AuthenticateState = 2
	ClientCredentialsError      AuthenticateState = 3
	DeviceFlowInitiated         AuthenticateState = 4
	DeviceFlowCodeReceived      AuthenticateState = 5
	DeviceFlowError             AuthenticateState = 6
	DeviceFlowBroswerWait       AuthenticateState = 7
	DeviceFlowBroswerCompleted  AuthenticateState = 8
	DeviceFlowBroswerError      AuthenticateState = 9
	DeviceFlowCompleted         AuthenticateState = 10
)

type LoginModel struct {
	styles                             *tui.Styles
	spinner                            spinner.Model
	width                              int
	height                             int
	optionsList                        list.Model
	isOptionChoosed                    bool
	optionChoosed                      string
	questionsForLoginAsMachine         []tui.Question
	currentLoginAsMachineQuestionIndex int
	loginAsMachineQuestionsDone        bool
	questionsForLoginAsUser            []tui.Question
	currentLoginAsUserQuestionIndex    int
	loginAsUserQuestionsDone           bool
	cli                                *core.CLI
	status                             AuthenticateState
	statusMessage                      string
	deviceFlowState                    auth.State
}

func (m LoginModel) runLoginAsMachine() error {

	err := core.RunLoginAsMachine(
		core.LoginInputs{
			ClientID:     m.questionsForLoginAsMachine[0].Answer,
			ClientSecret: m.questionsForLoginAsMachine[1].Answer,
			Tenant:       m.questionsForLoginAsMachine[2].Answer,
		}, m.cli)
	return err
}

func (m LoginModel) getDeviceCode() (auth.State, error) {

	return core.GetDeviceCode(m.cli)
}

func (m LoginModel) getAccessTokenFromDeviceCode(state auth.State) error {

	return core.GetAccessTokenFromDeviceCode(m.cli, state)
}

func NewLoginModel(cli *core.CLI) LoginModel {

	// Create a list of items for the user to choose from to authenticate
	items := []list.Item{
		tui.NewItem("As a machine", "Authenticates the IS CLI as a machine using client credentials"),
		tui.NewItem("As a user", "Authenticates the IS CLI as a user using personal credentials"),
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "How would you like to authenticate?"

	// Create a list of questions to ask the user when authenticating as a machine
	questionsForLoginAsMachine := []tui.Question{tui.NewQuestion("client id", "Client ID", tui.ShortQuestion), tui.NewQuestion("client secret", "Client Secret", tui.ShortSecreatQuestion), tui.NewQuestion("tenant", "Your tenant domain", tui.ShortQuestion)}

	// Create a list of questions to ask the user when authenticating as a user
	questionsForLoginAsUser := []tui.Question{tui.NewQuestion("client id", "Client ID", tui.ShortQuestion), tui.NewQuestion("tenant", "Your tenant domain", tui.ShortQuestion)}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return LoginModel{
		styles:                             tui.DefaultStyles(),
		spinner:                            s,
		optionsList:                        l,
		isOptionChoosed:                    false,
		questionsForLoginAsMachine:         questionsForLoginAsMachine,
		questionsForLoginAsUser:            questionsForLoginAsUser,
		currentLoginAsMachineQuestionIndex: 0,
		currentLoginAsUserQuestionIndex:    0,
		optionChoosed:                      "",
		cli:                                cli,
		status:                             NotStarted}

}

func (m LoginModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	currentLoginAsMachineQuestion := &m.questionsForLoginAsMachine[m.currentLoginAsMachineQuestionIndex]
	currentLoginAsUserQuestion := &m.questionsForLoginAsUser[m.currentLoginAsUserQuestionIndex]
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			if !m.isOptionChoosed {
				i, ok := m.optionsList.SelectedItem().(tui.Item)
				if ok {
					m.optionChoosed = i.Title()
					m.isOptionChoosed = true
				}
			} else {
				if m.optionChoosed == "As a user" {
					if m.status == DeviceFlowBroswerWait {
						m.status = DeviceFlowBroswerCompleted
						err := m.getAccessTokenFromDeviceCode(m.deviceFlowState)
						if err != nil {
							m.status = DeviceFlowError
							m.statusMessage = err.Error()
						} else {
							m.status = DeviceFlowCompleted
							return m, tea.Quit
						}
					} else {
						if m.currentLoginAsUserQuestionIndex == len(m.questionsForLoginAsUser)-1 {
							m.loginAsUserQuestionsDone = true
							currentLoginAsUserQuestion.Answer = currentLoginAsUserQuestion.Input.Value()
							m.status = DeviceFlowInitiated
							state, err := m.getDeviceCode()
							if err != nil {
								m.status = DeviceFlowError
								m.statusMessage = err.Error()
							} else {
								m.status = DeviceFlowCodeReceived
								if err = browser.OpenURL(state.VerificationURIComplete); err != nil {
									m.status = DeviceFlowBroswerError
								}
								m.deviceFlowState = state
								m.status = DeviceFlowBroswerWait
							}
						} else {
							m.NextLoginAsUserQuestion()
						}
						currentLoginAsUserQuestion.Answer = currentLoginAsUserQuestion.Input.Value()
						return m, currentLoginAsUserQuestion.Input.Blur
					}
				} else {
					if m.currentLoginAsMachineQuestionIndex == len(m.questionsForLoginAsMachine)-1 {
						m.loginAsMachineQuestionsDone = true
						currentLoginAsMachineQuestion.Answer = currentLoginAsMachineQuestion.Input.Value()
						m.status = ClientCredentialsInProgress
						err := m.runLoginAsMachine()
						if err != nil {
							m.statusMessage = err.Error()
							m.status = ClientCredentialsError
						} else {
							m.status = ClientCredentialsCompleted
						}
						return m, nil
					} else {
						m.NextLoginAsMachineQuestion()
					}
					currentLoginAsMachineQuestion.Answer = currentLoginAsMachineQuestion.Input.Value()
					return m, currentLoginAsMachineQuestion.Input.Blur
				}
			}

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := m.styles.List.GetFrameSize()
		m.optionsList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.optionsList, _ = m.optionsList.Update(msg)
	currentLoginAsMachineQuestion.Input, _ = currentLoginAsMachineQuestion.Input.Update(msg)
	currentLoginAsUserQuestion.Input, _ = currentLoginAsUserQuestion.Input.Update(msg)
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m LoginModel) View() string {
	if m.isOptionChoosed {
		if m.optionChoosed == "As a user" {
			current := m.questionsForLoginAsUser[m.currentLoginAsUserQuestionIndex]
			if !m.loginAsUserQuestionsDone {
				return current.Input.View()
			} else {
				if m.status == ClientCredentialsCompleted {
					return "Authenticated"
				} else if m.status == ClientCredentialsError {
					return "Error authenticating as a machine - " + m.statusMessage
				} else if m.status == ClientCredentialsInProgress {
					return "Authenticating as a user..."
				} else if m.status == DeviceFlowInitiated {
					return "Device flow initiated."
				} else if m.status == DeviceFlowBroswerWait {
					return fmt.Sprintf("\n\n   %s Waiting for the login to complete in the browser. Please press Enter after login completed!\n\n", m.spinner.View())
				} else if m.status == DeviceFlowBroswerCompleted {
					return "Device flow completed."
				} else if m.status == DeviceFlowBroswerError {
					return "Error opening browser. Please visit " + m.deviceFlowState.VerificationURIComplete + " to authenticate."
				} else if m.status == DeviceFlowCompleted {
					return "Successfully logged in"
				} else if m.status == DeviceFlowError {
					return "Error initiating device flow - " + m.statusMessage
				}
				return ""
			}
		} else {
			current := m.questionsForLoginAsMachine[m.currentLoginAsMachineQuestionIndex]
			if !m.loginAsMachineQuestionsDone {
				return current.Input.View()
			} else {
				if m.status == ClientCredentialsCompleted {
					return "Successfully authenticated as a machine"
				} else if m.status == ClientCredentialsError {
					return "Error authenticating as a machine - " + m.statusMessage
				} else if m.status == ClientCredentialsInProgress {
					return "Authenticating as a machine..."
				}
				return ""
			}
		}

	} else {
		return m.styles.List.Render(m.optionsList.View())
	}
}

func (m *LoginModel) NextLoginAsUserQuestion() {
	if m.currentLoginAsUserQuestionIndex < len(m.questionsForLoginAsMachine)-1 {
		m.currentLoginAsUserQuestionIndex++
	} else {
		m.currentLoginAsUserQuestionIndex = 0
	}
}

func (m *LoginModel) NextLoginAsMachineQuestion() {
	if m.currentLoginAsMachineQuestionIndex < len(m.questionsForLoginAsMachine)-1 {
		m.currentLoginAsMachineQuestionIndex++
	} else {
		m.currentLoginAsMachineQuestionIndex = 0
	}
}
