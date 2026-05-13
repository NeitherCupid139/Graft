package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveMigrationDirFindsServerRelativePathFromRepoRoot(t *testing.T) {
	root := t.TempDir()
	migrationDir := filepath.Join(root, "server", defaultMigrationDir)
	if err := os.MkdirAll(migrationDir, 0o755); err != nil {
		t.Fatalf("mkdir migration dir: %v", err)
	}

	resolved, err := resolveMigrationDir(root, defaultMigrationDir)
	if err != nil {
		t.Fatalf("resolve migration dir: %v", err)
	}

	if resolved != migrationDir {
		t.Fatalf("expected %s, got %s", migrationDir, resolved)
	}
}

func TestResolveMigrationDirFindsPathFromServerModuleRoot(t *testing.T) {
	root := t.TempDir()
	serverRoot := filepath.Join(root, "server")
	migrationDir := filepath.Join(serverRoot, defaultMigrationDir)
	if err := os.MkdirAll(migrationDir, 0o755); err != nil {
		t.Fatalf("mkdir migration dir: %v", err)
	}

	resolved, err := resolveMigrationDir(serverRoot, defaultMigrationDir)
	if err != nil {
		t.Fatalf("resolve migration dir: %v", err)
	}

	if resolved != migrationDir {
		t.Fatalf("expected %s, got %s", migrationDir, resolved)
	}
}

func TestResolveMigrationDirRejectsMissingPath(t *testing.T) {
	root := t.TempDir()

	_, err := resolveMigrationDir(root, defaultMigrationDir)
	if err == nil {
		t.Fatal("expected missing migration dir error")
	}
}
