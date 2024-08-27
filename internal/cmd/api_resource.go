package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shashimalcse/asgardeo-cli/internal/core"
	interactive "github.com/shashimalcse/asgardeo-cli/internal/interactive/api_resource"
	"github.com/spf13/cobra"
)

func apiResourceCmd(cli *core.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apis",
		Short: "Manage api resources",
	}

	cmd.AddCommand(listApiResourceCmd(cli))
	cmd.AddCommand(createAPIResourceCmd(cli))
	return cmd
}

func listApiResourceCmd(cli *core.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Short:   "List your api resources",
		Example: `asgardeo apis list
  asgardeo apis ls`,
		RunE: func(cmd *cobra.Command, args []string) error {

			m := interactive.NewApiResourceListModel(cli)
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

func createAPIResourceCmd(cli *core.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Args:    cobra.NoArgs,
		Short:   "Create an api resource",
		Example: `asgardeo apis create
  asgardeo apis c`,
		RunE: func(cmd *cobra.Command, args []string) error {

			m := interactive.NewAPIResourceCreateModel(cli)
			p := tea.NewProgram(m, tea.WithAltScreen())

			m1, err := p.Run()
			if err != nil {
				fmt.Println("Oh no:", err)
				os.Exit(1)
			}
			if m2, ok := m1.(interactive.APIResourceCreateModel); ok && m2.Value() != "" {
				fmt.Print(m2.Value())
			}
			return nil
		},
	}

	return cmd
}
