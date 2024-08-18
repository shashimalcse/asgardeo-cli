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
