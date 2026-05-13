package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

type devOptions struct {
	migrationDir string
}

var devMigrateRunner = func(cmd *cobra.Command, migrationDir string) error {
	return runMigrateUp(cmd, migrateUpOptions{migrationDir: migrationDir})
}

var devServeRunner = runServe

func newDevCommand() *cobra.Command {
	var opts devOptions

	command := &cobra.Command{
		Use:   "dev",
		Short: "Run migrations and start the Graft server for local development",
		Long: "graft dev is an explicit local development orchestration command. " +
			"It runs the migration CLI first and starts the server only after migrations succeed.",
		Example:      "  graft dev",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDev(cmd, args, opts)
		},
	}

	command.Flags().StringVar(&opts.migrationDir, "dir", defaultMigrationDir, "migration directory")
	return command
}

func runDev(cmd *cobra.Command, args []string, opts devOptions) error {
	if err := devMigrateRunner(cmd, opts.migrationDir); err != nil {
		return fmt.Errorf("run development migrations: %w", err)
	}

	if err := devServeRunner(cmd, args); err != nil {
		return fmt.Errorf("start development server: %w", err)
	}

	return nil
}
