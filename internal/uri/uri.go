package uri

import (
	"fmt"
	"regexp"
	"strings"
)

var validChars = regexp.MustCompile(`^[a-zA-Z0-9_.\-]+$`)

// Parsed represents a parsed kc:// URI.
type Parsed struct {
	Service string
	Key     string
}

// Parse parses a kc://service/key URI.
func Parse(raw string) (Parsed, error) {
	if !strings.HasPrefix(raw, "kc://") {
		return Parsed{}, fmt.Errorf("Invalid URI: %s (expected kc://service/key)", raw)
	}

	rest := strings.TrimPrefix(raw, "kc://")
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Parsed{}, fmt.Errorf("Invalid URI: %s (expected kc://service/key)", raw)
	}

	service := parts[0]
	key := parts[1]

	if strings.Contains(key, "/") {
		return Parsed{}, fmt.Errorf("Invalid URI: %s (expected kc://service/key)", raw)
	}

	if !validChars.MatchString(service) || !validChars.MatchString(key) {
		return Parsed{}, fmt.Errorf("Invalid characters in service/key (allowed: [a-zA-Z0-9_.-])")
	}

	return Parsed{Service: service, Key: key}, nil
}

// IsKCURI returns true if the string looks like a kc:// URI.
func IsKCURI(s string) bool {
	return strings.HasPrefix(s, "kc://")
}

// Format returns the kc://service/key string.
func Format(service, key string) string {
	return fmt.Sprintf("kc://%s/%s", service, key)
}
