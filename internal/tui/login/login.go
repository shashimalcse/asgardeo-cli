package login

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type Question struct {
	question string
	answer   string
	input    tui.Input
}

func newQuestion(q string) Question {
	return Question{question: q}
}

func NewShortQuestion(q string, p string) Question {
	question := newQuestion(q)
	model := tui.NewShortAnswerField(p)
	question.input = model
	return question
}

func NewLongQuestion(q string) Question {
	question := newQuestion(q)
	model := tui.NewLongAnswerField()
	question.input = model
	return question
}

type Model struct {
	styles                        *Styles
	width                         int
	height                        int
	list                          list.Model
	optionChoosed                 bool
	choice                        string
	asMachineQuestions            []Question
	currentAsMachineQuestionIndex int
	asMachineQuestionsDone        bool
}

func NewModel() Model {
	items := []list.Item{
		tui.NewItem("As a machine", "Authenticates the IS CLI as a machine using client credentials"),
		tui.NewItem("As a user", "Authenticates the IS CLI as a user using personal credentials"),
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "How would you like to authenticate?"

	questions := []Question{NewShortQuestion("client id", "Client ID"), NewShortQuestion("client secret", "Client Secret"), NewShortQuestion("tenant", "Your tenant domain")}

	return Model{list: l, optionChoosed: false, choice: "", asMachineQuestions: questions, currentAsMachineQuestionIndex: 0, styles: DefaultStyles()}

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	current := &m.asMachineQuestions[m.currentAsMachineQuestionIndex]
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			if !m.optionChoosed {
				i, ok := m.list.SelectedItem().(tui.Item)
				if ok {
					if i.Title() == "As a machine" {
						m.choice = "As a machine"

					} else if i.Title() == "As a user" {
						m.choice = "As a user"
					}
					m.optionChoosed = true
				}
			} else {
				if m.currentAsMachineQuestionIndex == len(m.asMachineQuestions)-1 {
					m.asMachineQuestionsDone = true
				}
				current.answer = current.input.Value()
				m.NextAsMachineQuestion()
				return m, current.input.Blur
			}

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := m.styles.List.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.optionChoosed {
		if m.asMachineQuestionsDone {
			return "Done"
		}
		current := m.asMachineQuestions[m.currentAsMachineQuestionIndex]
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Top,
			lipgloss.Left,
			lipgloss.JoinVertical(
				lipgloss.Left,
				current.question,
				m.styles.InputField.Render(current.input.View()),
			),
		)
	}
	return m.styles.List.Render(m.list.View())
}

func (m *Model) NextAsMachineQuestion() {
	if m.currentAsMachineQuestionIndex < len(m.asMachineQuestions)-1 {
		m.currentAsMachineQuestionIndex++
	} else {
		m.currentAsMachineQuestionIndex = 0
	}
}
