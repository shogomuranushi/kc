package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Entry represents a single line in a .env file.
type Entry struct {
	Key      string
	Value    string
	Raw      string // original line
	IsKCRef  bool   // true if value is a kc:// reference
	IsBlank  bool   // blank or comment line
	HasValue bool   // false if value is empty (KEY=)
}

// Parse reads and parses a .env file.
func Parse(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse .env: %s", err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			entries = append(entries, Entry{Raw: line, IsBlank: true})
			continue
		}

		// Strip optional "export " prefix
		kvLine := trimmed
		if strings.HasPrefix(kvLine, "export ") {
			kvLine = strings.TrimPrefix(kvLine, "export ")
			kvLine = strings.TrimSpace(kvLine)
		}

		eqIdx := strings.Index(kvLine, "=")
		if eqIdx < 0 {
			entries = append(entries, Entry{Raw: line, IsBlank: true})
			continue
		}

		key := strings.TrimSpace(kvLine[:eqIdx])
		val := strings.TrimSpace(kvLine[eqIdx+1:])

		// Remove surrounding quotes
		val = unquote(val)

		hasValue := val != ""
		isKCRef := strings.HasPrefix(val, "kc://")

		entries = append(entries, Entry{
			Key:     key,
			Value:   val,
			Raw:     line,
			IsKCRef: isKCRef,
			IsBlank: false,
			HasValue: hasValue,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Failed to parse .env: %s", err)
	}
	return entries, nil
}

// unquote removes surrounding single or double quotes.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// Write writes entries back to a .env file.
func Write(path string, entries []Entry) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i, e := range entries {
		if i > 0 {
			_, _ = w.WriteString("\n")
		}
		_, _ = w.WriteString(e.Raw)
	}
	_, _ = w.WriteString("\n")
	return w.Flush()
}

// FindEnvFile searches for a .env file starting from dir and going up.
func FindEnvFile(dir string) string {
	for {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err == nil {
			return path
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
