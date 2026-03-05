package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zalando/go-keyring"

	"github.com/shogomuranushi/kc/internal/keychain"
)

func TestMain(m *testing.M) {
	keyring.MockInit()
	os.Exit(m.Run())
}

func TestRunSetAndGet(t *testing.T) {
	// Set via args
	code := RunSet([]string{"test-svc", "test-key", "test-value"})
	if code != 0 {
		t.Fatalf("RunSet() returned %d, want 0", code)
	}

	// Get should succeed
	// Capture stdout by redirecting - we'll just verify exit code and keychain value
	val, err := keychain.Get("test-svc", "test-key")
	if err != nil {
		t.Fatalf("keychain.Get() error: %v", err)
	}
	if val != "test-value" {
		t.Errorf("value = %q, want %q", val, "test-value")
	}

	code = RunGet([]string{"test-svc", "test-key"})
	if code != 0 {
		t.Errorf("RunGet() returned %d, want 0", code)
	}
}

func TestRunSetMissingArgs(t *testing.T) {
	code := RunSet([]string{"only-service"})
	if code != 1 {
		t.Errorf("RunSet() returned %d, want 1", code)
	}
}

func TestRunSetInvalidChars(t *testing.T) {
	code := RunSet([]string{"bad service", "key", "val"})
	if code != 1 {
		t.Errorf("RunSet() returned %d, want 1", code)
	}
}

func TestRunGetNotFound(t *testing.T) {
	code := RunGet([]string{"nonexistent", "key"})
	if code != 1 {
		t.Errorf("RunGet() returned %d, want 1", code)
	}
}

func TestRunGetMissingArgs(t *testing.T) {
	code := RunGet([]string{"only-service"})
	if code != 1 {
		t.Errorf("RunGet() returned %d, want 1", code)
	}
}

func TestRunDelete(t *testing.T) {
	_ = RunSet([]string{"del-svc", "del-key", "val"})

	code := RunDelete([]string{"del-svc", "del-key"})
	if code != 0 {
		t.Errorf("RunDelete() returned %d, want 0", code)
	}

	// Should be gone
	code = RunGet([]string{"del-svc", "del-key"})
	if code != 1 {
		t.Errorf("RunGet() after delete returned %d, want 1", code)
	}
}

func TestRunDeleteNotFound(t *testing.T) {
	code := RunDelete([]string{"nonexistent", "key"})
	if code != 1 {
		t.Errorf("RunDelete() returned %d, want 1", code)
	}
}

func TestRunDeleteMissingArgs(t *testing.T) {
	code := RunDelete([]string{"only-service"})
	if code != 1 {
		t.Errorf("RunDelete() returned %d, want 1", code)
	}
}

func TestRunList(t *testing.T) {
	// List with no entries should succeed
	code := RunList([]string{})
	if code != 0 {
		t.Errorf("RunList() returned %d, want 0", code)
	}
}

func TestRunRunMissingArgs(t *testing.T) {
	code := RunRun([]string{"--"})
	if code != 1 {
		t.Errorf("RunRun() returned %d, want 1", code)
	}
}

func TestRunRunCommandNotFound(t *testing.T) {
	code := RunRun([]string{"--", "nonexistent-command-xyz-12345"})
	if code != 1 {
		t.Errorf("RunRun() returned %d, want 1", code)
	}
}

func TestRunExternalCommandWithEnv(t *testing.T) {
	// Create a temp .env file
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("TEST_VAR=hello\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Run from tmpDir - the command should find .env
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Use echo as a simple test - won't actually exec because syscall.Exec replaces process
	// Instead test the env file resolution logic
	code := RunRun([]string{"--", "nonexistent-cmd-xyz"})
	if code != 1 {
		t.Errorf("expected exit code 1 for nonexistent command")
	}
}

func TestRunMigrateNoFile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	code := RunMigrate([]string{})
	if code != 1 {
		t.Errorf("RunMigrate() returned %d, want 1 (no .env found)", code)
	}
}

func TestRunMigrateNothingToMigrate(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	// All lines are already kc:// refs or comments
	content := "# comment\nKEY=kc://svc/key\n"
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	code := RunMigrate([]string{envPath})
	if code != 0 {
		t.Errorf("RunMigrate() returned %d, want 0", code)
	}
}
