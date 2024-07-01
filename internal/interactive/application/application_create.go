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

// Application Create Interactive

type ApplicationCreateState int

const (
	ApplicationCreateInitiated ApplicationCreateState = iota
	ApplicationCreateTypeSelected
	ApplicationCreateQuestionsCompleted
	ApplicationCreateCreatingInProgress
	ApplicationCreateCreatingCompleted
	ApplicationCreateError
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
	confirmationQuestion           tui.Question
	currentSinglePageQuestionIndex int
	applicationType                ApplicationType
	output                         string
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

	confirmationQuestion := tui.NewQuestion("Are you sure you want to create the application? (y/n)", "(y/n)", tui.ShortQuestion)

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
		confirmationQuestion:           confirmationQuestion,
		currentSinglePageQuestionIndex: 0,
	}

}

func (m ApplicationCreateModel) createApplications() error {

	if m.applicationType == SinglePage {
		application := management.Application{
			Name:       m.questionsForSinglePage[0].Answer,
			TemplateID: "6a90e4b0-fbff-42d7-bfde-1efd98f07cd7",
			AdvancedConfig: management.AdvancedConfigurations{
				DiscoverableByEndUsers: false,
				SkipLoginConsent:       true,
				SkipLogoutConsent:      true,
			},
			AssociatedRoles: management.AssociatedRoles{
				AllowedAudience: "APPLICATION",
				Roles:           []management.AssociatedRole{},
			},
			AuthenticationSeq: management.AuthenticationSequence{
				Type: "DEFAULT",
				Steps: []management.Step{{
					ID: 1,
					Options: []management.Options{
						{IDP: "LOCAL", Authenticator: "basic"},
					},
				},
				},
			},
			ClaimConfiguration: management.ClaimConfiguration{
				Dialect: "LOCAL",
				RequestedClaims: []interface{}{
					map[string]interface{}{
						"claim": map[string]interface{}{"uri": "http://wso2.org/claims/username"},
					},
				},
			},
			InboundProtocolConfiguration: management.InboundProtocolConfiguration{
				OIDC: management.OIDC{
					AccessToken: management.AccessToken{
						ApplicationAccessTokenExpiryInSeconds: 3600,
						BindingType:                           "sso-session",
						RevokeTokensWhenIDPSessionTerminated:  true,
						Type:                                  "Default",
						UserAccessTokenExpiryInSeconds:        3600,
						ValidateTokenBinding:                  false,
					},
					AllowedOrigins: []string{m.questionsForSinglePage[1].Answer},
					CallbackURLs:   []string{m.questionsForSinglePage[1].Answer},
					GrantTypes:     []string{"authorization_code", "refresh_token"},
					PKCE: management.PKCE{
						Mandatory:                      true,
						SupportPlainTransformAlgorithm: false,
					},
					PublicClient: true,
					RefreshToken: management.RefreshToken{
						ExpiryInSeconds:   86400,
						RenewRefreshToken: true,
					},
				},
			},
		}
		_, err := m.cli.API.Application.Create(context.Background(), application)
		return err
	}
	return nil
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
						m.confirmationQuestion.Input.SetValue("")
					} else {
						m.NextSinglePageQuestion()
					}
					return m, currentSinglePageQuestion.Input.Blur
				}
			case ApplicationCreateQuestionsCompleted:
				m.confirmationQuestion.Answer = m.confirmationQuestion.Input.Value()
				if (m.confirmationQuestion.Answer == "y") || (m.confirmationQuestion.Answer == "Y" || m.confirmationQuestion.Answer == "") {
					m.state = ApplicationCreateCreatingInProgress
					err := m.createApplications()
					if err != nil {
						m.state = ApplicationCreateError
						m.stateError = err
						m.output = "Error creating application!"
					} else {
						m.state = ApplicationCreateCreatingCompleted
						m.output = "Application created successfully!"
					}
				} else {
					m.output = "Application creation cancelled."
					return m, tea.Quit
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
	m.confirmationQuestion.Input, _ = m.confirmationQuestion.Input.Update(msg)
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
		return previousQAs + m.confirmationQuestion.Question + "\n" + m.confirmationQuestion.Input.View()
	case ApplicationCreateCreatingInProgress:
		return fmt.Sprintf("\n\n   %s Creating application...!\n\n", m.spinner.View())
	case ApplicationCreateCreatingCompleted:
		return "Application created successfully!"
	case ApplicationCreateError:
		return fmt.Sprint(m.stateError.Error())
	}

	return ""
}

func (m ApplicationCreateModel) Value() string {
	return fmt.Sprint(m.output)
}

func (m *ApplicationCreateModel) NextSinglePageQuestion() {
	if m.currentSinglePageQuestionIndex < len(m.questionsForSinglePage)-1 {
		m.currentSinglePageQuestionIndex++
	} else {
		m.currentSinglePageQuestionIndex = 0
	}
}
