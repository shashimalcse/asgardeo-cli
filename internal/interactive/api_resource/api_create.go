package interactive

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shashimalcse/asgardeo-cli/internal/core"
	"github.com/shashimalcse/asgardeo-cli/internal/models"
	"github.com/shashimalcse/asgardeo-cli/internal/tui"
)

type APIResourceCreateState int

const (
	StateInitiated APIResourceCreateState = iota
	StateAddingScopes
	StateConfirmation
	StateCreatingInProgress
	StateCreatingCompleted
	StateCreatingError
)

type APIResourceCreateModel struct {
	styles               *tui.Styles
	spinner              spinner.Model
	width, height        int
	cli                  *core.CLI
	state                APIResourceCreateState
	stateError           error
	questions            []tui.Question
	currentQuestionIndex int
	scopes               []models.Scope
	scopeInput           textinput.Model
	output               string
}

func NewAPIResourceCreateModel(cli *core.CLI) *APIResourceCreateModel {
	m := &APIResourceCreateModel{
		styles:    tui.DefaultStyles(),
		spinner:   newSpinner(),
		cli:       cli,
		state:     StateInitiated,
		questions: initQuestions(),
		scopes:    []models.Scope{},
	}
	m.scopeInput = textinput.New()
	m.scopeInput.Placeholder = "Enter a scope (or leave empty to finish)"
	m.scopeInput.Focus()
	return m
}

func initQuestions() []tui.Question {
	questions := []tui.Question{
		tui.NewQuestion("Identifier", "Identifier", tui.ShortQuestion),
		tui.NewQuestion("Display Name", "Display Name", tui.ShortQuestion),
		tui.NewQuestion("Confirm", "Do you want to create this API Resource? (Y/n)", tui.ShortQuestion),
	}
	return questions
}

func (m *APIResourceCreateModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *APIResourceCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m.handleKeyEnter()
		}
	case tea.WindowSizeMsg:
		return m.handleWindowResize(msg)
	}

	var cmd tea.Cmd

	switch m.state {
	case StateInitiated, StateConfirmation:
		m.questions[m.currentQuestionIndex].Input, cmd = m.questions[m.currentQuestionIndex].Input.Update(msg)
	case StateAddingScopes:
		m.scopeInput, cmd = m.scopeInput.Update(msg)
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m *APIResourceCreateModel) handleKeyEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case StateInitiated:
		currentQuestion := &m.questions[m.currentQuestionIndex]
		currentQuestion.Answer = currentQuestion.Input.Value()
		if m.currentQuestionIndex == len(m.questions)-2 {
			m.state = StateAddingScopes
			m.scopeInput.SetValue("")
		} else {
			m.NextQuestion()
		}
		return m, currentQuestion.Input.Blur
	case StateAddingScopes:
		scopeName := strings.TrimSpace(m.scopeInput.Value())
		if scopeName != "" {
			m.scopes = append(m.scopes, models.Scope{
				Description: "",
				DisplayName: scopeName,
				Name:        scopeName,
			})
			m.scopeInput.SetValue("")
			return m, nil
		}
		m.state = StateConfirmation
		m.currentQuestionIndex = len(m.questions) - 1
		return m, nil
	case StateConfirmation:
		confirmation := strings.ToLower(m.questions[m.currentQuestionIndex].Input.Value())
		if confirmation == "y" || confirmation == "Y" || confirmation == "" {
			m.state = StateCreatingInProgress
			err := m.createAPIResources()
			if err != nil {
				m.state = StateCreatingError
				m.stateError = err
				m.output = "Error creating APIResource!"
			} else {
				m.state = StateCreatingCompleted
				m.output = "APIResource created successfully!"
			}
		} else {
			m.output = "APIResource creation cancelled."
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *APIResourceCreateModel) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width, m.height = msg.Width, msg.Height
	_, _ = m.styles.List.GetFrameSize()
	return m, nil
}

func (m *APIResourceCreateModel) View() string {
	switch m.state {
	case StateInitiated:
		return m.renderQuestions()
	case StateAddingScopes:
		return m.renderScopeInput()
	case StateConfirmation:
		return m.renderConfirmation()
	case StateCreatingInProgress:
		return fmt.Sprintf("\n\n   %s Creating APIResource...\n\n", m.spinner.View())
	case StateCreatingCompleted:
		return "APIResource created successfully!"
	case StateCreatingError:
		return fmt.Sprintf("Error creating APIResource: %v", m.stateError)
	}
	return ""
}

func (m *APIResourceCreateModel) renderQuestions() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Creating a new API Resource \n\n"))
	for i, q := range m.questions[:m.currentQuestionIndex] {
		sb.WriteString(fmt.Sprintf("%s: %s\n", q.Question, q.Answer))
		if i == len(m.questions)-1 {
			sb.WriteString("\n")
		}
	}
	sb.WriteString(m.questions[m.currentQuestionIndex].Input.View())
	return sb.String()
}

func (m *APIResourceCreateModel) renderScopeInput() string {
	var sb strings.Builder
	sb.WriteString("Enter scopes for the API Resource (press Enter with an empty input to finish):\n\n")
	if len(m.scopes) > 0 {
		sb.WriteString("Scopes:\n")
		for i, scope := range m.scopes {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, scope.Name))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	sb.WriteString(m.scopeInput.View())
	return sb.String()
}

func (m *APIResourceCreateModel) renderConfirmation() string {
	var sb strings.Builder
	sb.WriteString("API Resource Details:\n\n")
	for _, q := range m.questions[:len(m.questions)-1] {
		sb.WriteString(fmt.Sprintf("%s: %s\n", q.Question, q.Answer))
	}
	sb.WriteString("\nScopes:\n")
	if len(m.scopes) > 0 {
		for i, scope := range m.scopes {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, scope.Name))
		}
	} else {
		sb.WriteString("No scopes defined.\n")
	}
	sb.WriteString("\n")
	sb.WriteString(m.questions[m.currentQuestionIndex].Input.View())
	return sb.String()
}

func (m *APIResourceCreateModel) Value() string {
	return fmt.Sprint(m.output)
}

func (m *APIResourceCreateModel) NextQuestion() {
	if m.currentQuestionIndex < len(m.questions)-1 {
		m.currentQuestionIndex++
	} else {
		m.currentQuestionIndex = 0
	}
}

func (m *APIResourceCreateModel) createAPIResources() error {
	payload := map[string]interface{}{
		"identifier":            m.questions[0].Answer,
		"name":                  m.questions[1].Answer,
		"requiresAuthorization": true,
		"scopes":                m.scopes,
	}
	err := m.cli.API.APIResource.Create(context.Background(), payload)
	if err != nil {
		return err
	}
	return nil
}
