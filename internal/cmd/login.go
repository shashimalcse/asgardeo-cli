package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shashimalcse/is-cli/internal/core"
	"github.com/shashimalcse/is-cli/internal/interactive"
	"github.com/shashimalcse/is-cli/internal/models"
	"github.com/spf13/cobra"
)

func loginCmd(cli *core.CLI) *cobra.Command {
	var inputs core.LoginInputs
	var verbose bool

	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Authenticate the IS CLI",
		Example: "is login",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine if we should use interactive mode
			if !inputs.IsLoggingInAsAMachine() {
				result := runInteractiveLogin(cli)
				if result.IsError {
					return fmt.Errorf(result.Message)
				} else {
					fmt.Println(result.Message)
					return nil
				}
			}
			return runMachineLogin(cli, inputs, verbose)
		},
	}

	cmd.Flags().StringVar(&inputs.ClientID, "client-id", "", "Client ID")
	cmd.Flags().StringVar(&inputs.ClientSecret, "client-secret", "", "Client Secret")
	cmd.Flags().StringVar(&inputs.Tenant, "tenant", "", "Tenant")
	cmd.MarkFlagsRequiredTogether("client-id", "client-secret", "tenant")

	return cmd
}

func runInteractiveLogin(cli *core.CLI) models.OutputResult {
	m := interactive.NewLoginModel(cli)
	p := tea.NewProgram(m, tea.WithAltScreen())

	m1, err := p.Run()
	if err != nil {
		return models.OutputResult{
			Message: fmt.Sprintf("Interactive login failed: %v", err),
			IsError: true,
		}
	}
	if m2, ok := m1.(interactive.LoginModel); ok {
		return m2.GetOutputValue()
	}
	return models.OutputResult{}
}

func runMachineLogin(cli *core.CLI, inputs core.LoginInputs, verbose bool) error {
	if verbose {
		fmt.Println("Attempting machine login...")
	}

	if err := validateMachineLoginInputs(inputs); err != nil {
		return err
	}

	if err := core.RunLoginAsMachine(inputs, cli); err != nil {
		return fmt.Errorf("failed to login as machine: %w", err)
	}

	if verbose {
		fmt.Println("Machine login successful")
	}
	return nil
}

func validateMachineLoginInputs(inputs core.LoginInputs) error {
	if inputs.ClientID == "" || inputs.ClientSecret == "" || inputs.Tenant == "" {
		return fmt.Errorf("client-id, client-secret, and tenant are required for machine login")
	}
	return nil
}
