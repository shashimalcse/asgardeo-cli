package interactive

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shashimalcse/asgardeo-cli/internal/core"
	"github.com/shashimalcse/asgardeo-cli/internal/tui"
)

type OrganizationCreateState int

const (
	StateInitiated OrganizationCreateState = iota
	StateConfirmation
	StateCreatingInProgress
	StateCreatingCompleted
	StateCreatingError
)

type OrganizationCreateModel struct {
	styles               *tui.Styles
	spinner              spinner.Model
	width, height        int
	cli                  *core.CLI
	state                OrganizationCreateState
	stateError           error
	questions            []tui.Question
	currentQuestionIndex int
	output               string
}

func NewOrganizationCreateModel(cli *core.CLI) *OrganizationCreateModel {
	m := &OrganizationCreateModel{
		styles:    tui.DefaultStyles(),
		spinner:   newSpinner(),
		cli:       cli,
		state:     StateInitiated,
		questions: initQuestions(),
	}
	return m
}

func initQuestions() []tui.Question {
	questions := []tui.Question{
		tui.NewQuestion("Name", "Organization Name", tui.ShortQuestion),
		tui.NewQuestion("Description", "Description (optional)", tui.ShortQuestion),
		tui.NewQuestion("Confirm", "Do you want to create this organization? (Y/n)", tui.ShortQuestion),
	}
	return questions
}

func (m *OrganizationCreateModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *OrganizationCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case string:
		m.state = StateCreatingCompleted
		m.output = msg
		return m, nil
	case error:
		m.state = StateCreatingError
		m.stateError = msg
		return m, nil
	}

	var cmd tea.Cmd

	switch m.state {
	case StateInitiated, StateConfirmation:
		m.questions[m.currentQuestionIndex].Input, cmd = m.questions[m.currentQuestionIndex].Input.Update(msg)
	}

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	return m, tea.Batch(cmd, spinnerCmd)
}

func (m *OrganizationCreateModel) handleKeyEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case StateInitiated:
		currentQuestion := &m.questions[m.currentQuestionIndex]
		currentQuestion.Answer = currentQuestion.Input.Value()
		// Move to next question or confirmation state
		if m.currentQuestionIndex < len(m.questions)-2 {
			m.currentQuestionIndex++
		} else {
			// We've finished all non-confirmation questions
			m.state = StateConfirmation
			m.currentQuestionIndex = len(m.questions) - 1
		}
	case StateConfirmation:
		answer := strings.ToLower(m.questions[m.currentQuestionIndex].Input.Value())
		if answer == "y" || answer == "" {
			m.state = StateCreatingInProgress
			return m, m.createOrganization
		}
		return m, tea.Quit
	case StateCreatingCompleted, StateCreatingError:
		return m, tea.Quit
	}
	return m, nil
}

func (m *OrganizationCreateModel) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	return m, nil
}

func (m *OrganizationCreateModel) createOrganization() tea.Msg {
	ctx := context.Background()

	organization := make(map[string]interface{})
	organization["name"] = m.questions[0].Answer

	if m.questions[1].Answer != "" {
		organization["description"] = m.questions[1].Answer
	}

	err := m.cli.API.Organization.Create(ctx, organization)
	if err != nil {
		return err
	}

	return "Organization created successfully"
}

func (m *OrganizationCreateModel) View() string {
	switch m.state {
	case StateInitiated:
		return m.renderQuestions()
	case StateConfirmation:
		return m.renderConfirmation()
	case StateCreatingInProgress:
		return fmt.Sprintf("\n\n   %s Creating organization...\n\n", m.spinner.View())
	case StateCreatingCompleted:
		return m.output
	case StateCreatingError:
		return fmt.Sprintf("Error: %v", m.stateError)
	}
	return ""
}

func (m *OrganizationCreateModel) renderQuestions() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Creating a new Organization \n\n"))
	for _, q := range m.questions[:m.currentQuestionIndex] {
		sb.WriteString(fmt.Sprintf("%s: %s\n", q.Question, q.Answer))
	}
	sb.WriteString(m.questions[m.currentQuestionIndex].Input.View())
	return sb.String()
}

func (m *OrganizationCreateModel) renderConfirmation() string {
	var sb strings.Builder
	sb.WriteString("Organization Details:\n\n")
	sb.WriteString(fmt.Sprintf("Name: %s\n", m.questions[0].Answer))
	if m.questions[1].Answer != "" {
		sb.WriteString(fmt.Sprintf("Description: %s\n", m.questions[1].Answer))
	}
	sb.WriteString("\n")
	sb.WriteString(m.questions[m.currentQuestionIndex].Input.View())
	return sb.String()
}

func (m *OrganizationCreateModel) Value() string {
	return m.output
}
