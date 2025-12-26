package cmd

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shashimalcse/asgardeo-cli/internal/core"
	interactive "github.com/shashimalcse/asgardeo-cli/internal/interactive/organization"
	"github.com/shashimalcse/asgardeo-cli/internal/models"
	"github.com/spf13/cobra"
)

func organizationsCmd(cli *core.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "Manage organizations",
	}

	cmd.AddCommand(listOrganizationsCmd(cli))
	cmd.AddCommand(createOrganizationCmd(cli))
	cmd.AddCommand(updateOrganizationCmd(cli))
	cmd.AddCommand(updateOrganizationMetadataCmd(cli))
	cmd.AddCommand(deleteOrganizationCmd(cli))
	return cmd
}

func listOrganizationsCmd(cli *core.CLI) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Short:   "List your organizations",
		Example: `asgardeo orgs list
  asgardeo orgs ls
  asgardeo orgs list --filter "name eq MyOrg"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			m := interactive.NewOrganizationListModel(cli)
			p := tea.NewProgram(m, tea.WithAltScreen())

			if _, err := p.Run(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&filter, "filter", "", "Filter organizations (e.g., 'name eq MyOrg')")
	return cmd
}

func createOrganizationCmd(cli *core.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Args:    cobra.NoArgs,
		Short:   "Create an organization",
		Example: `asgardeo orgs create
  asgardeo orgs c`,
		RunE: func(cmd *cobra.Command, args []string) error {
			m := interactive.NewOrganizationCreateModel(cli)
			p := tea.NewProgram(m, tea.WithAltScreen())
			m1, err := p.Run()
			if err != nil {
				fmt.Println("Oh no:", err)
				os.Exit(1)
			}
			if m2, ok := m1.(*interactive.OrganizationCreateModel); ok && m2.Value() != "" {
				fmt.Print(m2.Value())
			}
			return nil
		},
	}
	return cmd
}

type OrganizationUpdateInputs struct {
	OrganizationId string
	Name           string
	Description    string
}

func updateOrganizationCmd(cli *core.CLI) *cobra.Command {
	var inputs OrganizationUpdateInputs
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Update organization information",
		Example: `asgardeo orgs update <org-id> --name "New Name" --description "New Description"
  asgardeo orgs u <org-id> --name "New Name"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				inputs.OrganizationId = args[0]
			}
			if inputs.OrganizationId == "" {
				return fmt.Errorf("organization ID is required")
			}
			if inputs.Name == "" && inputs.Description == "" {
				return fmt.Errorf("at least one field (name or description) must be provided")
			}

			updates := make(map[string]interface{})
			if inputs.Name != "" {
				updates["name"] = inputs.Name
			}
			if inputs.Description != "" {
				updates["description"] = inputs.Description
			}

			fmt.Printf("Updating organization with ID: %s\n", inputs.OrganizationId)
			err := cli.API.Organization.Update(context.Background(), inputs.OrganizationId, updates)
			if err != nil {
				return err
			}
			fmt.Println("Organization updated successfully")
			return nil
		},
	}
	cmd.Flags().StringVar(&inputs.OrganizationId, "org-id", "", "Organization ID")
	cmd.Flags().StringVar(&inputs.Name, "name", "", "New organization name")
	cmd.Flags().StringVar(&inputs.Description, "description", "", "New organization description")
	return cmd
}

type OrganizationMetadataInputs struct {
	OrganizationId string
	Operation      string
	Path           string
	Value          string
}

func updateOrganizationMetadataCmd(cli *core.CLI) *cobra.Command {
	var inputs OrganizationMetadataInputs
	cmd := &cobra.Command{
		Use:     "update-metadata",
		Aliases: []string{"um"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Update organization metadata",
		Example: `asgardeo orgs update-metadata <org-id> --operation replace --path "/attributes/0/value" --value "new-value"
  asgardeo orgs um <org-id> --operation add --path "/attributes" --value "key:value"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				inputs.OrganizationId = args[0]
			}
			if inputs.OrganizationId == "" {
				return fmt.Errorf("organization ID is required")
			}
			if inputs.Operation == "" || inputs.Path == "" || inputs.Value == "" {
				return fmt.Errorf("operation, path, and value are required")
			}

			metadata := []models.OrganizationPatch{
				{
					Operation: inputs.Operation,
					Path:      inputs.Path,
					Value:     inputs.Value,
				},
			}

			fmt.Printf("Updating metadata for organization with ID: %s\n", inputs.OrganizationId)
			err := cli.API.Organization.UpdateMetadata(context.Background(), inputs.OrganizationId, metadata)
			if err != nil {
				return err
			}
			fmt.Println("Organization metadata updated successfully")
			return nil
		},
	}
	cmd.Flags().StringVar(&inputs.OrganizationId, "org-id", "", "Organization ID")
	cmd.Flags().StringVar(&inputs.Operation, "operation", "", "Operation (add, replace, remove)")
	cmd.Flags().StringVar(&inputs.Path, "path", "", "JSON path to update")
	cmd.Flags().StringVar(&inputs.Value, "value", "", "Value to set")
	return cmd
}

type OrganizationDeleteInputs struct {
	OrganizationId string
}

func deleteOrganizationCmd(cli *core.CLI) *cobra.Command {
	var inputs OrganizationDeleteInputs
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Delete an organization",
		Example: `asgardeo orgs delete <org-id>
  asgardeo orgs rm <org-id>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				inputs.OrganizationId = args[0]
			}
			if inputs.OrganizationId == "" {
				return fmt.Errorf("organization ID is required")
			}
			fmt.Printf("Deleting organization with ID: %s\n", inputs.OrganizationId)
			err := cli.API.Organization.Delete(context.Background(), inputs.OrganizationId)
			if err != nil {
				return err
			}
			fmt.Println("Organization deleted successfully")
			return nil
		},
	}
	cmd.Flags().StringVar(&inputs.OrganizationId, "org-id", "", "Organization ID")
	return cmd
}
