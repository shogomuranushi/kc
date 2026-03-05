package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shogomuranushi/kc/internal/dotenv"
	"github.com/shogomuranushi/kc/internal/keychain"
	"github.com/shogomuranushi/kc/internal/uri"
)

func RunMigrate(args []string) int {
	var envFilePath string

	if len(args) > 0 {
		envFilePath = args[0]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
			return 1
		}
		envFilePath = dotenv.FindEnvFile(cwd)
		if envFilePath == "" {
			fmt.Fprintln(os.Stderr, "[ERROR] No .env file found")
			return 1
		}
	}

	entries, err := dotenv.Parse(envFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		return 1
	}

	// Find plaintext entries to migrate
	var targets []int
	for i, e := range entries {
		if e.IsBlank || !e.HasValue || e.IsKCRef {
			continue
		}
		targets = append(targets, i)
	}

	if len(targets) == 0 {
		fmt.Fprintln(os.Stderr, "No plaintext values to migrate")
		return 0
	}

	absPath, _ := os.Getwd()
	if !strings.HasPrefix(envFilePath, "/") {
		envFilePath = absPath + "/" + envFilePath
	}
	fmt.Fprintf(os.Stderr, "\nFound %d plaintext values in %s\n\n", len(targets), envFilePath)

	reader := bufio.NewReader(os.Stdin)

	for idx, i := range targets {
		e := entries[i]
		// Mask value for display
		display := maskValue(e.Value)
		fmt.Fprintf(os.Stderr, "[%d/%d] %s=%s\n", idx+1, len(targets), e.Key, display)

		// Suggest service name from key
		suggestedService := suggestService(e.Key)
		suggestedKey := suggestKey(e.Key)

		svc := prompt(reader, "  Service name", suggestedService)
		key := prompt(reader, "  Key name", suggestedKey)

		if !validNameChars.MatchString(svc) || !validNameChars.MatchString(key) {
			fmt.Fprintln(os.Stderr, "  [ERROR] Invalid characters in service/key (allowed: [a-zA-Z0-9_.-])")
			continue
		}

		if err := keychain.Set(svc, key, e.Value); err != nil {
			fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
			continue
		}

		kcURI := uri.Format(svc, key)
		fmt.Fprintf(os.Stderr, "  ✓ Saved %s\n\n", kcURI)

		// Update entry
		entries[i].Raw = fmt.Sprintf("%s=%s", e.Key, kcURI)
		entries[i].Value = kcURI
		entries[i].IsKCRef = true
	}

	// Write updated .env
	if err := dotenv.Write(envFilePath, entries); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to write .env: %s\n", err)
		return 1
	}

	fmt.Fprintln(os.Stderr, ".env updated.")
	return 0
}

func prompt(reader *bufio.Reader, label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Fprintf(os.Stderr, "%s [%s]: ", label, defaultVal)
	} else {
		fmt.Fprintf(os.Stderr, "%s: ", label)
	}
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal
	}
	return line
}

func maskValue(val string) string {
	if len(val) <= 8 {
		return "****"
	}
	return val[:4] + "..." + val[len(val)-4:]
}

func suggestService(envKey string) string {
	key := strings.ToLower(envKey)
	// Common patterns
	for _, prefix := range []string{"github", "aws", "stripe", "anthropic", "openai", "slack", "docker"} {
		if strings.Contains(key, prefix) {
			return prefix
		}
	}
	// Use first part before underscore
	parts := strings.SplitN(key, "_", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func suggestKey(envKey string) string {
	key := strings.ToLower(envKey)
	// Remove common service prefixes
	for _, prefix := range []string{"github_", "aws_", "stripe_", "anthropic_", "openai_", "slack_", "docker_"} {
		if strings.HasPrefix(key, prefix) {
			return strings.TrimPrefix(key, prefix)
		}
	}
	return key
}

