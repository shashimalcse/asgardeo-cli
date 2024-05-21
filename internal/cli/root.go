package cli

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

const rootShort = "Build, manage and test your Identity Server integrations from the command line."

func Execute() {
	cli := &cli{}

	cobra.EnableCommandSorting = false
	rootCmd := buildRootCmd(cli)
	addSubCommands(rootCmd, cli)

	cancelCtx := contextWithCancel()
	if err := rootCmd.ExecuteContext(cancelCtx); err != nil {

		os.Exit(1) // nolint:gocritic
	}
}

func buildRootCmd(cli *cli) *cobra.Command {

	rootCommand := &cobra.Command{
		Use:           "is",
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:         rootShort,
		Long:          rootShort,
		Version:       "v0.0.1",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return rootCommand
}

func addSubCommands(rootCmd *cobra.Command, cli *cli) {

	rootCmd.AddCommand(loginCmd(cli))
}

func contextWithCancel() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		defer cancel()
		os.Exit(0)
	}()

	return ctx
}
