package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shashimalcse/is-cli/internal/tui/login"
	"github.com/spf13/cobra"
)

func loginCmd(cli *cli) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to the Identity Server",
		RunE: func(cmd *cobra.Command, args []string) error {

			m := login.NewModel()

			p := tea.NewProgram(m, tea.WithAltScreen())

			if _, err := p.Run(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}
			return nil
		},
	}

	return cmd
}
