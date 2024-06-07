package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shashimalcse/is-cli/internal/core"
	"github.com/shashimalcse/is-cli/internal/interactive"
	"github.com/spf13/cobra"
)

func applicationsCmd(cli *core.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "applications",
		Short: "Manage applications",
	}

	cmd.AddCommand(listApplicationsCmd(cli))

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

			m := interactive.NewApplicationModel(cli)
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
