package interactive

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shashimalcse/is-cli/internal/auth"
	"github.com/shashimalcse/is-cli/internal/core"
	"github.com/shashimalcse/is-cli/internal/tui"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
	List        lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	s.List = lipgloss.NewStyle().Margin(1, 2)
	return s
}

type AuthenticateState int

const (
	NotStarted          AuthenticateState = 0
	InProgress          AuthenticateState = 1
	Completed           AuthenticateState = 2
	Error               AuthenticateState = 3
	DeviceFlowInitiated AuthenticateState = 4
)

type Model struct {
	styles                             *Styles
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

func (m Model) runLoginAsMachine() error {

	err := core.RunLoginAsMachine(
		core.LoginInputs{
			ClientID:     m.questionsForLoginAsMachine[0].Answer,
			ClientSecret: m.questionsForLoginAsMachine[1].Answer,
			Tenant:       m.questionsForLoginAsMachine[2].Answer,
		}, m.cli)
	return err
}

func (m Model) getDeviceCode() (auth.State, error) {

	return core.GetDeviceCode(m.cli)
}

func NewModel(cli *core.CLI) Model {

	// Create a list of items for the user to choose from to authenticate
	items := []list.Item{
		tui.NewItem("As a machine", "Authenticates the IS CLI as a machine using client credentials"),
		tui.NewItem("As a user", "Authenticates the IS CLI as a user using personal credentials"),
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "How would you like to authenticate?"

	// Create a list of questions to ask the user when authenticating as a machine
	questionsForLoginAsMachine := []tui.Question{tui.NewQuestion("client id", "Client ID", tui.ShortQuestion), tui.NewQuestion("client secret", "Client Secret", tui.ShortQuestion), tui.NewQuestion("tenant", "Your tenant domain", tui.ShortQuestion)}

	// Create a list of questions to ask the user when authenticating as a user
	questionsForLoginAsUser := []tui.Question{tui.NewQuestion("client id", "Client ID", tui.ShortQuestion), tui.NewQuestion("tenant", "Your tenant domain", tui.ShortQuestion)}

	return Model{
		styles:                             DefaultStyles(),
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

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
					if m.currentLoginAsUserQuestionIndex == len(m.questionsForLoginAsUser)-1 {
						m.loginAsUserQuestionsDone = true
						currentLoginAsUserQuestion.Answer = currentLoginAsUserQuestion.Input.Value()
						state, err := m.getDeviceCode()
						if err != nil {
							m.statusMessage = err.Error()
						} else {
							m.deviceFlowState = state

						}
					} else {
						m.NextLoginAsUserQuestion()
					}
					currentLoginAsUserQuestion.Answer = currentLoginAsUserQuestion.Input.Value()
					return m, currentLoginAsUserQuestion.Input.Blur
				} else {
					if m.currentLoginAsMachineQuestionIndex == len(m.questionsForLoginAsMachine)-1 {
						m.loginAsMachineQuestionsDone = true
						currentLoginAsMachineQuestion.Answer = currentLoginAsMachineQuestion.Input.Value()
						m.status = InProgress
						err := m.runLoginAsMachine()
						if err != nil {
							m.statusMessage = err.Error()
							m.status = Error
						} else {
							m.status = Completed
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
	m.optionsList, cmd = m.optionsList.Update(msg)
	currentLoginAsMachineQuestion.Input, cmd = currentLoginAsMachineQuestion.Input.Update(msg)
	currentLoginAsUserQuestion.Input, cmd = currentLoginAsUserQuestion.Input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.isOptionChoosed {
		if m.optionChoosed == "As a user" {
			current := m.questionsForLoginAsUser[m.currentLoginAsUserQuestionIndex]
			if !m.loginAsUserQuestionsDone {
				return current.Input.View()
			} else {
				if m.status == Completed {
					return "Authenticated"
				} else if m.status == Error {
					return "Error authenticating as a machine - " + m.statusMessage
				} else if m.status == InProgress {
					return "Authenticating as a user..."
				}
				return ""
			}
		} else {
			current := m.questionsForLoginAsMachine[m.currentLoginAsMachineQuestionIndex]
			if !m.loginAsMachineQuestionsDone {
				return current.Input.View()
			} else {
				if m.status == Completed {
					return "Authenticated as a machine"
				} else if m.status == Error {
					return "Error authenticating as a machine - " + m.statusMessage
				} else if m.status == InProgress {
					return "Authenticating as a machine..."
				}
				return ""
			}
		}

	} else {
		return m.styles.List.Render(m.optionsList.View())
	}
}

func (m *Model) NextLoginAsUserQuestion() {
	if m.currentLoginAsUserQuestionIndex < len(m.questionsForLoginAsMachine)-1 {
		m.currentLoginAsUserQuestionIndex++
	} else {
		m.currentLoginAsUserQuestionIndex = 0
	}
}

func (m *Model) NextLoginAsMachineQuestion() {
	if m.currentLoginAsMachineQuestionIndex < len(m.questionsForLoginAsMachine)-1 {
		m.currentLoginAsMachineQuestionIndex++
	} else {
		m.currentLoginAsMachineQuestionIndex = 0
	}
}
