package cmd

import (
	"context"
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
	cmd.AddCommand(deleteAPIResourceCmd(cli))
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
			if m2, ok := m1.(*interactive.APIResourceCreateModel); ok && m2.Value() != "" {
				fmt.Print(m2.Value())
			}
			return nil
		},
	}

	return cmd
}

type APIResourceDeleteInputs struct {
	ApiId string
}

func deleteAPIResourceCmd(cli *core.CLI) *cobra.Command {
	var inputs APIResourceDeleteInputs
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Delete api resource",
		Example: `asgardeo apis delete
                  asgardeo apis rm`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				inputs.ApiId = args[0]
			}
			if inputs.ApiId == "" {
				return fmt.Errorf("api resource ID is required")
			}
			fmt.Printf("Deleting api resource with ID: %s\n", inputs.ApiId)
			err := cli.API.APIResource.Delete(context.Background(), inputs.ApiId)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&inputs.ApiId, "api-id", "", "API Resource ID")
	return cmd
}
