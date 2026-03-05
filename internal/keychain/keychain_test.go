package keychain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestMain(m *testing.M) {
	keyring.MockInit()
	os.Exit(m.Run())
}

func setupRegistry(t *testing.T) {
	t.Helper()
	// Use temp dir for registry
	tmpDir := t.TempDir()
	origFunc := registryPathFunc
	registryPathFunc = func() string {
		return filepath.Join(tmpDir, "registry.json")
	}
	t.Cleanup(func() {
		registryPathFunc = origFunc
	})
}

func TestSetAndGet(t *testing.T) {
	setupRegistry(t)

	err := Set("github", "token", "ghp_xxxx")
	if err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	val, err := Get("github", "token")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if val != "ghp_xxxx" {
		t.Errorf("Get() = %q, want %q", val, "ghp_xxxx")
	}
}

func TestGetNotFound(t *testing.T) {
	setupRegistry(t)

	_, err := Get("nonexistent", "key")
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
}

func TestSetOverwrite(t *testing.T) {
	setupRegistry(t)

	_ = Set("svc", "key", "old")
	_ = Set("svc", "key", "new")

	val, _ := Get("svc", "key")
	if val != "new" {
		t.Errorf("Get() = %q, want %q", val, "new")
	}
}

func TestDelete(t *testing.T) {
	setupRegistry(t)

	_ = Set("svc", "key", "val")
	err := Delete("svc", "key")
	if err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	_, err = Get("svc", "key")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestDeleteNotFound(t *testing.T) {
	setupRegistry(t)

	err := Delete("nonexistent", "key")
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
}

func TestList(t *testing.T) {
	setupRegistry(t)

	_ = Set("aws", "access_key_id", "AKIA")
	_ = Set("aws", "secret_access_key", "xxxx")
	_ = Set("github", "token", "ghp")

	entries, err := List("")
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("List() returned %d entries, want 3", len(entries))
	}

	// Filter by service
	awsEntries, err := List("aws")
	if err != nil {
		t.Fatalf("List(aws) error: %v", err)
	}
	if len(awsEntries) != 2 {
		t.Fatalf("List(aws) returned %d entries, want 2", len(awsEntries))
	}
}

func TestListEmpty(t *testing.T) {
	setupRegistry(t)

	entries, err := List("")
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("List() returned %d entries, want 0", len(entries))
	}
}
