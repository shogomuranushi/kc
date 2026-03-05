package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	content := `# Comment line
PLAIN_VALUE=hello
SECRET=kc://service/key
QUOTED="kc://service/key2"
SINGLE_QUOTED='plain_value'
export EXPORTED=exported_value
EMPTY_VALUE=
PORT=3000

# Another comment
`
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	entries, err := Parse(envPath)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	// Expected: comment, PLAIN_VALUE, SECRET, QUOTED, SINGLE_QUOTED, export EXPORTED, EMPTY_VALUE, PORT, blank, comment, (trailing blank from last \n)
	// Let's check specific entries by finding them
	var kvEntries []Entry
	for _, e := range entries {
		if !e.IsBlank {
			kvEntries = append(kvEntries, e)
		}
	}

	if len(kvEntries) != 7 {
		t.Fatalf("expected 7 key-value entries, got %d", len(kvEntries))
	}

	tests := []struct {
		idx      int
		key      string
		value    string
		isKCRef  bool
		hasValue bool
	}{
		{0, "PLAIN_VALUE", "hello", false, true},
		{1, "SECRET", "kc://service/key", true, true},
		{2, "QUOTED", "kc://service/key2", true, true},
		{3, "SINGLE_QUOTED", "plain_value", false, true},
		{4, "EXPORTED", "exported_value", false, true},
		{5, "EMPTY_VALUE", "", false, false},
		{6, "PORT", "3000", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			e := kvEntries[tt.idx]
			if e.Key != tt.key {
				t.Errorf("key = %q, want %q", e.Key, tt.key)
			}
			if e.Value != tt.value {
				t.Errorf("value = %q, want %q", e.Value, tt.value)
			}
			if e.IsKCRef != tt.isKCRef {
				t.Errorf("isKCRef = %v, want %v", e.IsKCRef, tt.isKCRef)
			}
			if e.HasValue != tt.hasValue {
				t.Errorf("hasValue = %v, want %v", e.HasValue, tt.hasValue)
			}
		})
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestFindEnvFile(t *testing.T) {
	// Create nested directories with .env at top level
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("KEY=val\n"), 0644); err != nil {
		t.Fatal(err)
	}

	subDir := filepath.Join(tmpDir, "a", "b", "c")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Should find .env from nested subdirectory
	found := FindEnvFile(subDir)
	if found != envPath {
		t.Errorf("FindEnvFile() = %q, want %q", found, envPath)
	}
}

func TestFindEnvFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	// No .env file exists
	found := FindEnvFile(tmpDir)
	if found != "" {
		t.Errorf("FindEnvFile() = %q, want empty string", found)
	}
}

func TestFindEnvFileCurrentDir(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("KEY=val\n"), 0644); err != nil {
		t.Fatal(err)
	}

	found := FindEnvFile(tmpDir)
	if found != envPath {
		t.Errorf("FindEnvFile() = %q, want %q", found, envPath)
	}
}

func TestWrite(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	entries := []Entry{
		{Raw: "# comment", IsBlank: true},
		{Key: "KEY", Value: "kc://svc/key", Raw: "KEY=kc://svc/key", IsKCRef: true, HasValue: true},
		{Key: "PLAIN", Value: "hello", Raw: "PLAIN=hello", HasValue: true},
	}

	if err := Write(envPath, entries); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatal(err)
	}

	expected := "# comment\nKEY=kc://svc/key\nPLAIN=hello\n"
	if string(data) != expected {
		t.Errorf("Write() output = %q, want %q", string(data), expected)
	}
}
