package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"graft/server/internal/config"
)

const defaultMigrationDir = "internal/ent/migrate/migrations"

var migrateLookPath = exec.LookPath
var migrateCommandContext = exec.CommandContext
var migrateGetwd = os.Getwd
var migrateStdin io.Reader = os.Stdin

type migrateUpOptions struct {
	migrationDir string
	workingDir   string
}

func newMigrateCommand() *cobra.Command {
	var migrationDir string

	command := &cobra.Command{
		Use:   "migrate",
		Short: "Run explicit database migration commands",
	}
	command.PersistentFlags().StringVar(&migrationDir, "dir", defaultMigrationDir, "migration directory")

	command.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Apply pending Atlas versioned migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrateUp(cmd, migrateUpOptions{migrationDir: migrationDir})
		},
	})

	return command
}

func runMigrateUp(cmd *cobra.Command, opts migrateUpOptions) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	workingDir := opts.workingDir
	if strings.TrimSpace(workingDir) == "" {
		workingDir, err = migrateGetwd()
		if err != nil {
			return fmt.Errorf("resolve working directory: %w", err)
		}
	}

	absDir, err := resolveMigrationDir(workingDir, opts.migrationDir)
	if err != nil {
		return fmt.Errorf("resolve migration dir: %w", err)
	}

	atlasPath, err := findAtlasCLI()
	if err != nil {
		return err
	}

	command := migrateCommandContext(
		cmd.Context(),
		atlasPath,
		"migrate",
		"apply",
		"--dir", "file://"+filepath.ToSlash(absDir),
		"--url", cfg.Database.URL,
	)
	command.Stdout = cmd.OutOrStdout()
	command.Stderr = cmd.ErrOrStderr()
	command.Stdin = migrateStdin

	if err := command.Run(); err != nil {
		return fmt.Errorf("apply atlas migrations: %w", err)
	}

	return nil
}

func findAtlasCLI() (string, error) {
	atlasPath, err := migrateLookPath("atlas")
	if err == nil {
		return atlasPath, nil
	}

	return "", fmt.Errorf(
		"atlas CLI is required for `graft migrate up` and `graft dev`; install Atlas first, or run `graft serve` only after the database schema is already up to date: %w",
		err,
	)
}

func resolveMigrationDir(baseDir string, migrationDir string) (string, error) {
	if strings.TrimSpace(migrationDir) == "" {
		return "", fmt.Errorf("migration dir is required")
	}

	searchDirs := []string{migrationDir}
	if migrationDir == defaultMigrationDir {
		searchDirs = append(searchDirs, filepath.Join("server", migrationDir))
	}

	current := baseDir
	for {
		for _, relativeDir := range searchDirs {
			candidate := filepath.Join(current, relativeDir)
			info, err := os.Stat(candidate)
			if err == nil {
				if !info.IsDir() {
					return "", fmt.Errorf("migration dir %s is not a directory", candidate)
				}

				absDir, err := filepath.Abs(candidate)
				if err != nil {
					return "", fmt.Errorf("resolve migration dir %s: %w", candidate, err)
				}

				return absDir, nil
			}
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf("cannot find migration dir %q from %s", migrationDir, baseDir)
}
