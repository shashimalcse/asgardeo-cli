package cli

import (
	"context"
	"os"
	"os/signal"

	"github.com/shashimalcse/is-cli/internal/core"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const rootShort = `

 _____    _            _   _ _           _____                          
|_   _|  | |          | | (_) |         /  ___|                         
  | |  __| | ___ _ __ | |_ _| |_ _   _  \ ` + "`" + `--.  ___ _ ____   _____ _ __ 
  | | / _` + "`" + ` |/ _ \ '_ \| __| | __| | | |  ` + "`" + `--. \/ _ \ '__\ \ / / _ \ '__|
 _| || (_| |  __/ | | | |_| | |_| |_| | /\__/ /  __/ |   \ V /  __/ |   
|___/ \__,_|\___|_| |_|\__|_|\__|\__, | \____/ \___|_|    \_/ \___|_|   
                                 __/  |                                 
                                |____/                                  

Build, manage and test your Identity Server integrations from the command line.								
`

func Execute() {
	cli := &core.CLI{
		Logger: configLogger(),
	}

	cobra.EnableCommandSorting = false
	rootCmd := buildRootCmd(cli)
	addSubCommands(rootCmd, cli)

	cancelCtx := contextWithCancel()
	if err := rootCmd.ExecuteContext(cancelCtx); err != nil {

		os.Exit(1) // nolint:gocritic
	}
}

func buildRootCmd(cli *core.CLI) *cobra.Command {

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

func addSubCommands(rootCmd *cobra.Command, cli *core.CLI) {

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

func configLogger() zap.Logger {

	config := zap.NewProductionConfig()

	// Set the log file path
	logFilePath := "is-cli.log"

	// Create a file to write logs to
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Set up the zap logger to write to the specified file
	config.OutputPaths = []string{logFilePath}
	config.ErrorOutputPaths = []string{logFilePath}

	// Build the logger
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // Flushes buffer, if any

	return *logger
}
