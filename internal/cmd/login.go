package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shashimalcse/is-cli/internal/core"
	"github.com/shashimalcse/is-cli/internal/interactive"
	"github.com/spf13/cobra"
)

func loginCmd(cli *core.CLI) *cobra.Command {

	var inputs core.LoginInputs
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Authenticate the IS CLI",
		Example: "is login",
		RunE: func(cmd *cobra.Command, args []string) error {

			shouldPrompt := !inputs.IsLoggingInAsAMachine()
			if shouldPrompt {

				m := interactive.NewModel(cli)
				p := tea.NewProgram(m, tea.WithAltScreen())

				if _, err := p.Run(); err != nil {
					fmt.Println("Error running program:", err)
					os.Exit(1)
				}
			} else {
				if err := core.RunLoginAsMachine(core.LoginInputs{ClientID: inputs.ClientID, ClientSecret: inputs.ClientSecret, Tenant: inputs.Tenant}, cli); err != nil {
					return err
				}
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
