package keychain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zalando/go-keyring"
)

const servicePrefix = "kc-cli:"

// Set stores a secret in the keychain.
func Set(service, key, value string) error {
	svc := servicePrefix + service
	err := keyring.Set(svc, key, value)
	if err != nil {
		return fmt.Errorf("Keychain error: %w", err)
	}
	registryAdd(service, key)
	return nil
}

// Get retrieves a secret from the keychain.
func Get(service, key string) (string, error) {
	svc := servicePrefix + service
	val, err := keyring.Get(svc, key)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", fmt.Errorf("Not found: kc://%s/%s", service, key)
		}
		return "", fmt.Errorf("Keychain error: %w", err)
	}
	return val, nil
}

// Delete removes a secret from the keychain.
func Delete(service, key string) error {
	svc := servicePrefix + service
	err := keyring.Delete(svc, key)
	if err != nil {
		if err == keyring.ErrNotFound {
			return fmt.Errorf("Not found: kc://%s/%s", service, key)
		}
		return fmt.Errorf("Keychain error: %w", err)
	}
	registryRemove(service, key)
	return nil
}

// List returns kc:// URIs for stored secrets, optionally filtered by service.
func List(service string) ([]string, error) {
	entries := registryLoad()
	var result []string
	for _, e := range entries {
		parts := strings.SplitN(e, "/", 2)
		if len(parts) != 2 {
			continue
		}
		if service != "" && parts[0] != service {
			continue
		}
		result = append(result, fmt.Sprintf("kc://%s", e))
	}
	sort.Strings(result)
	return result, nil
}

// --- registry: tracks service/key pairs in a JSON file ---

var registryPathFunc = defaultRegistryPath

func defaultRegistryPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "kc-cli", "registry.json")
}

func registryLoad() []string {
	data, err := os.ReadFile(registryPathFunc())
	if err != nil {
		return nil
	}
	var entries []string
	_ = json.Unmarshal(data, &entries)
	return entries
}

func registrySave(entries []string) {
	p := registryPathFunc()
	_ = os.MkdirAll(filepath.Dir(p), 0700)
	data, _ := json.Marshal(entries)
	_ = os.WriteFile(p, data, 0600)
}

func registryAdd(service, key string) {
	entry := service + "/" + key
	entries := registryLoad()
	for _, e := range entries {
		if e == entry {
			return
		}
	}
	entries = append(entries, entry)
	sort.Strings(entries)
	registrySave(entries)
}

func registryRemove(service, key string) {
	entry := service + "/" + key
	entries := registryLoad()
	var filtered []string
	for _, e := range entries {
		if e != entry {
			filtered = append(filtered, e)
		}
	}
	registrySave(filtered)
}
