package cli

import (
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shashimalcse/is-cli/internal/auth"
	"github.com/shashimalcse/is-cli/internal/tui/login"
	"github.com/spf13/cobra"
)

type LoginInputs struct {
	ClientID     string
	ClientSecret string
	Tenant       string
}

func (i *LoginInputs) isLoggingInAsAMachine() bool {
	return i.ClientID != "" || i.ClientSecret != "" || i.Tenant != ""
}

func loginCmd(cli *cli) *cobra.Command {

	var inputs LoginInputs
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Authenticate the IS CLI",
		Example: "is login",
		RunE: func(cmd *cobra.Command, args []string) error {

			selectedLoginType := ""
			shouldPrompt := !inputs.isLoggingInAsAMachine()
			if shouldPrompt {

				m := login.NewModel(&selectedLoginType)
				p := tea.NewProgram(m, tea.WithAltScreen())

				if _, err := p.Run(); err != nil {
					fmt.Println("Error running program:", err)
					os.Exit(1)
				}
			}

			if selectedLoginType == "As a user" {
				fmt.Println("Logging in as a user")
			} else {
				result, err := auth.GetAccessTokenFromClientCreds(http.DefaultClient, auth.ClientCredentials{ClientID: inputs.ClientID, ClientSecret: inputs.ClientSecret, Tenant: inputs.Tenant})
				if err != nil {
					fmt.Println("Error logging in:", err)
				}
				fmt.Println("Access Token:", result.AccessToken)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&inputs.ClientID, "client-id", "", "", "Client ID")
	cmd.Flags().StringVarP(&inputs.ClientSecret, "client-secret", "", "", "Client Secret")
	cmd.Flags().StringVarP(&inputs.Tenant, "tenant", "", "", "Tenant")
	cmd.MarkFlagsRequiredTogether("client-id", "client-secret", "tenant")
	return cmd
}
