package cmd

import (
	"fmt"

	"github.com/shashimalcse/asgardeo-cli/internal/core"
	"github.com/shashimalcse/asgardeo-cli/internal/keyring"
	"github.com/spf13/cobra"
)

func logoutCmd(cli *core.CLI) *cobra.Command {

	var tenant string
	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "Logout the Asgardeo CLI",
		Example: "asgardeo logout <tenant>",
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := cli.Config.RemoveTenant(tenant); err != nil {
				return fmt.Errorf("failed to log out from the tenant %q: %w", tenant, err)
			}
			if err := keyring.DeleteSecretsForTenant(tenant); err != nil {
				return fmt.Errorf("failed to delete tenant secrets: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&tenant, "tenant", "", "", "tenant")
	cmd.MarkFlagsOneRequired("tenant")
	return cmd
}
