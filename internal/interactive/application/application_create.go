package interactive

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shashimalcse/asgardeo-cli/internal/core"
	"github.com/shashimalcse/asgardeo-cli/internal/tui"
)

var (
	OIDC = "OIDC"
	SAML = "SAML"
)

type ApplicationCreateState int

const (
	StateInitiated ApplicationCreateState = iota
	StateTypeSelected
	StateQuestionsCompleted
	StateCreatingInProgress
	StateCreatingCompleted
	StateCreatingError
)

type ApplicationType string

const (
	SinglePage       ApplicationType = "Single-Page Application"
	Traditional_OIDC ApplicationType = "Traditional Web Application OIDC"
	Traditional_SAML ApplicationType = "Traditional Web Application SAML"
	Mobile           ApplicationType = "Mobile Application"
	Standard         ApplicationType = "Standard-Based Application"
	M2M              ApplicationType = "M2M Application"
)

type ApplicationCreateModel struct {
	styles               *tui.Styles
	spinner              spinner.Model
	width, height        int
	cli                  *core.CLI
	state                ApplicationCreateState
	stateError           error
	applicationTypes     list.Model
	questions            []tui.Question
	currentQuestionIndex int
	applicationType      ApplicationType
	output               string
}

func NewApplicationCreateModel(cli *core.CLI) *ApplicationCreateModel {
	return &ApplicationCreateModel{
		styles:           tui.DefaultStyles(),
		spinner:          newSpinner(),
		cli:              cli,
		state:            StateInitiated,
		applicationTypes: newApplicationTypesList(),
	}
}

func newApplicationTypesList() list.Model {
	items := []list.Item{
		tui.NewItemWithKey("single_page", string(SinglePage), "A web application that runs application logic in the browser."),
		tui.NewItemWithKey("traditional", string(Traditional_OIDC), "A web application that runs application logic on the server."),
		tui.NewItemWithKey("mobile", string(Mobile), "Applications developed to target mobile devices."),
		tui.NewItemWithKey("standard", string(Standard), "Applications built using standard protocols."),
		tui.NewItemWithKey("m2m", string(M2M), "Applications tailored for Machine to Machine communication."),
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select application template to create application"
	return l
}

func (m *ApplicationCreateModel) initQuestions() []tui.Question {
	if m.applicationType == SinglePage {
		m.questions = []tui.Question{
			tui.NewQuestion("Name", "Name", tui.ShortQuestion),
			tui.NewQuestion("Authorized redirect URL", "Authorized redirect URL", tui.ShortQuestion),
			tui.NewQuestion("Are you sure you want to create the application? (y/n)", "Are you sure you want to create the application? (Y/n)", tui.ShortQuestion),
		}
	} else if m.applicationType == Traditional_OIDC {
		m.questions = []tui.Question{
			tui.NewQuestion("Name", "Name", tui.ShortQuestion),
			tui.NewQuestion("Protocol (OIDC/SAML)", "Protocol (OIDC/SAML) default : OIDC", tui.ShortQuestion),
			tui.NewQuestion("Authorized redirect URL", "Authorized redirect URL", tui.ShortQuestion),
			tui.NewQuestion("Are you sure you want to create the application? (Y/n)", "Are you sure you want to create the application? (Y/n)", tui.ShortQuestion),
		}
	} else if m.applicationType == Traditional_SAML {
		m.questions = []tui.Question{
			tui.NewQuestion("Name", "Name", tui.ShortQuestion),
			tui.NewQuestion("Protocol (OIDC/SAML)", "Protocol (OIDC/SAML) default : SAML", tui.ShortQuestion),
			tui.NewQuestion("Issuer", "Issuer", tui.ShortQuestion),
			tui.NewQuestion("Assertion consumer service URLs", "Assertion consumer service URLs", tui.ShortQuestion),
			tui.NewQuestion("Are you sure you want to create the application? (Y/n)", "Are you sure you want to create the application? (Y/n)", tui.ShortQuestion),
		}
	}
	return nil
}

func (m ApplicationCreateModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m ApplicationCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
	if m.state == StateInitiated {
		m.applicationTypes, _ = m.applicationTypes.Update(msg)
	}
	if m.state == StateTypeSelected || m.state == StateQuestionsCompleted {
		m.questions[m.currentQuestionIndex].Input, _ = m.questions[m.currentQuestionIndex].Input.Update(msg)
	}
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m ApplicationCreateModel) handleKeyEnter(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateInitiated:
		i, ok := m.applicationTypes.SelectedItem().(tui.Item)
		if ok {
			m.applicationType = ApplicationType(i.Title())
			m.initQuestions()
			m.state = StateTypeSelected
		}
	case StateTypeSelected:
		currentQuestion := &m.questions[m.currentQuestionIndex]
		currentQuestion.Answer = currentQuestion.Input.Value()
		if m.currentQuestionIndex == len(m.questions)-2 {
			m.state = StateQuestionsCompleted
			m.NextQuestion()
			m.questions[m.currentQuestionIndex].Input.SetValue("")
		} else {
			if m.questions[m.currentQuestionIndex].Question == "Protocol (OIDC/SAML)" {
				protocol := strings.TrimSpace(m.questions[m.currentQuestionIndex].Answer)
				protocol = strings.ToUpper(protocol)
				switch protocol {
				case OIDC:
					m.applicationType = Traditional_OIDC
				case SAML:
					m.applicationType = Traditional_SAML
				case "":
					m.applicationType = Traditional_OIDC
					m.questions[m.currentQuestionIndex].Answer = protocol
				default:
					m.output = "Invalid protocol. Please enter OIDC or SAML"
					return m, tea.Quit
				}
				m.initQuestions()
			}
			m.NextQuestion()
		}
		return m, currentQuestion.Input.Blur
	case StateQuestionsCompleted:
		confirmation := strings.ToLower(m.questions[m.currentQuestionIndex].Input.Value())
		if (confirmation == "y") || (confirmation == "Y" || confirmation == "") {
			m.state = StateCreatingInProgress
			err := m.createApplications()
			if err != nil {
				m.state = StateCreatingError
				m.stateError = err
				m.output = "Error creating application!"
			} else {
				m.state = StateCreatingCompleted
				m.output = "Application created successfully!"
			}
		} else {
			m.output = "Application creation cancelled."
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ApplicationCreateModel) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width, m.height = msg.Width, msg.Height
	h, v := m.styles.List.GetFrameSize()
	m.applicationTypes.SetSize(msg.Width-h, msg.Height-v)
	return m, nil
}

func (m ApplicationCreateModel) View() string {
	switch m.state {
	case StateInitiated:
		return m.styles.List.Render(m.applicationTypes.View())
	case StateTypeSelected, StateQuestionsCompleted:
		return m.renderQuestions()
	case StateCreatingInProgress:
		return fmt.Sprintf("\n\n   %s Creating application...\n\n", m.spinner.View())
	case StateCreatingCompleted:
		return "Application created successfully!"
	case StateCreatingError:
		return fmt.Sprintf("Error creating application: %v", m.stateError)
	}
	return ""
}

func (m *ApplicationCreateModel) renderQuestions() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Creating a new %s\n\n", m.applicationType))
	for i, q := range m.questions[:m.currentQuestionIndex] {
		sb.WriteString(fmt.Sprintf("%s: %s\n", q.Question, q.Answer))
		if i == len(m.questions)-1 {
			sb.WriteString("\n")
		}
	}
	sb.WriteString(m.questions[m.currentQuestionIndex].Input.View())
	return sb.String()
}

func (m ApplicationCreateModel) Value() string {
	return fmt.Sprint(m.output)
}

func (m *ApplicationCreateModel) NextQuestion() {
	if m.currentQuestionIndex < len(m.questions)-1 {
		m.currentQuestionIndex++
	} else {
		m.currentQuestionIndex = 0
	}
}

func (m ApplicationCreateModel) createApplications() error {

	application := map[string]interface{}{
		"name": m.questions[0].Answer,
		"advancedConfigurations": map[string]interface{}{
			"discoverableByEndUsers": false,
			"skipLogoutConsent":      true,
			"skipLoginConsent":       true,
		},
		"authenticationSequence": map[string]interface{}{
			"type": "DEFAULT",
			"steps": []interface{}{
				map[string]interface{}{
					"id": 1,
					"options": []interface{}{
						map[string]interface{}{
							"idp":           "LOCAL",
							"authenticator": "basic",
						},
					},
				},
			},
		},
		"associatedRoles": map[string]interface{}{
			"allowedAudience": "APPLICATION",
			"roles":           []string{},
		},
	}

	if m.applicationType == SinglePage {
		application["templateId"] = "6a90e4b0-fbff-42d7-bfde-1efd98f07cd7"
		application["inboundProtocolConfiguration"] = map[string]interface{}{
			"oidc": map[string]interface{}{
				"accessToken": map[string]interface{}{
					"applicationAccessTokenExpiryInSeconds": 3600,
					"bindingType":                           "sso-session",
					"revokeTokensWhenIDPSessionTerminated":  true,
					"type":                                  "Default",
					"userAccessTokenExpiryInSeconds":        3600,
					"validateTokenBinding":                  false,
				},
				"allowedOrigins": []string{m.questions[1].Answer},
				"callbackURLs":   []string{m.questions[1].Answer},
				"grantTypes":     []string{"authorization_code", "refresh_token"},
				"pkce": map[string]interface{}{
					"mandatory":                      true,
					"supportPlainTransformAlgorithm": false,
				},
				"publicClient": true,
				"refreshToken": map[string]interface{}{
					"expiryInSeconds":   86400,
					"renewRefreshToken": true,
				},
			},
		}
		application["claimConfiguration"] = map[string]interface{}{
			"dialect": "LOCAL",
			"requestedClaims": []interface{}{
				map[string]interface{}{
					"claim": map[string]interface{}{
						"uri": "http://wso2.org/claims/username",
					},
				},
			},
		}
	} else if m.applicationType == Traditional_OIDC {
		if m.questions[1].Answer == OIDC {
			application["templateId"] = "b9c5e11e-fc78-484b-9bec-015d247561b8"
			application["inboundProtocolConfiguration"] = map[string]interface{}{
				"oidc": map[string]interface{}{
					"allowedOrigins": []string{},
					"callbackURLs":   []string{m.questions[2].Answer},
					"grantTypes":     []string{"authorization_code"},
					"publicClient":   false,
					"refreshToken": map[string]interface{}{
						"expiryInSeconds": 86400,
					},
				},
			}
			application["claimConfiguration"] = map[string]interface{}{
				"dialect": "LOCAL",
				"requestedClaims": []interface{}{
					map[string]interface{}{
						"claim": map[string]interface{}{
							"uri": "http://wso2.org/claims/username",
						},
					},
				},
			}
		}
		if m.questions[1].Answer == SAML {
			application["templateId"] = "776a73da-fd8e-490b-84ff-93009f8ede85"
			application["inboundProtocolConfiguration"] = map[string]interface{}{
				"saml": map[string]interface{}{
					"manualConfiguration": map[string]interface{}{
						"issuer":                "https://localhost:9443/oauth2/token",
						"assertionConsumerUrls": []string{},
						"attributeProfile": map[string]interface{}{
							"alwaysIncludeAttributesInResponse": true,
							"enabled":                           true,
						},
						"singleLogoutProfile": map[string]interface{}{
							"enabled":      true,
							"logoutMethod": "BACKCHANNEL",
							"idpInitiatedSingleLogout": map[string]interface{}{
								"enabled": false,
							},
						},
					},
				},
			}
		}
	}
	err := m.cli.API.Application.Create(context.Background(), application)
	if err != nil {
		return err
	}
	return nil
}
