package cmd

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shashimalcse/is-cli/internal/core"
	interactive "github.com/shashimalcse/is-cli/internal/interactive/application"
	"github.com/spf13/cobra"
)

func applicationsCmd(cli *core.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "applications",
		Short: "Manage applications",
	}

	cmd.AddCommand(listApplicationsCmd(cli))
	cmd.AddCommand(createApplicationsCmd(cli))
	cmd.AddCommand(deleteApplicationsCmd(cli))
	return cmd
}

func listApplicationsCmd(cli *core.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Short:   "List your applications",
		Example: `is applications list
  is applications ls`,
		RunE: func(cmd *cobra.Command, args []string) error {

			m := interactive.NewApplicationListModel(cli)
			p := tea.NewProgram(m, tea.WithAltScreen())

			if _, err := p.Run(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			} else {

			}

			return nil
		},
	}

	return cmd
}

func createApplicationsCmd(cli *core.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Args:    cobra.NoArgs,
		Short:   "Create an application",
		Example: `is applications create
  is applications c`,
		RunE: func(cmd *cobra.Command, args []string) error {

			m := interactive.NewApplicationCreateModel(cli)
			p := tea.NewProgram(m, tea.WithAltScreen())

			m1, err := p.Run()
			if err != nil {
				fmt.Println("Oh no:", err)
				os.Exit(1)
			}
			if m2, ok := m1.(interactive.ApplicationCreateModel); ok && m2.Value() != "" {
				fmt.Print(m2.Value())
			}
			return nil
		},
	}

	return cmd
}

type ApplicationDeleteInputs struct {
	ApplicationId string
}

func deleteApplicationsCmd(cli *core.CLI) *cobra.Command {
	var inputs ApplicationDeleteInputs
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Delete applications",
		Example: `is applications delete
  is applications rm`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				inputs.ApplicationId = args[0]
			}
			if inputs.ApplicationId == "" {
				return fmt.Errorf("application ID is required")
			}
			fmt.Printf("Deleting application with ID: %s\n", inputs.ApplicationId)
			err := cli.API.Application.Delete(context.Background(), inputs.ApplicationId)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&inputs.ApplicationId, "app-id", "", "Application ID")
	return cmd
}
