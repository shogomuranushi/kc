package cmd

import (
	"fmt"
	"os"

	"github.com/shogomuranushi/kc/internal/keychain"
	"github.com/shogomuranushi/kc/internal/uri"
)

func RunDelete(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "[ERROR] Usage: kc delete <service> <key>")
		return 1
	}

	service := args[0]
	key := args[1]

	if !validNameChars.MatchString(service) || !validNameChars.MatchString(key) {
		fmt.Fprintln(os.Stderr, "[ERROR] Invalid characters in service/key (allowed: [a-zA-Z0-9_.-])")
		return 1
	}

	if err := keychain.Delete(service, key); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		return 1
	}

	fmt.Fprintf(os.Stderr, "Deleted: %s\n", uri.Format(service, key))
	return 0
}
