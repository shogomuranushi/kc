package cmd

import (
	"fmt"
	"os"

	"github.com/shogomuranushi/kc/internal/keychain"
)

func RunGet(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "[ERROR] Usage: kc get <service> <key>")
		return 1
	}

	service := args[0]
	key := args[1]

	if !validNameChars.MatchString(service) || !validNameChars.MatchString(key) {
		fmt.Fprintln(os.Stderr, "[ERROR] Invalid characters in service/key (allowed: [a-zA-Z0-9_.-])")
		return 1
	}

	val, err := keychain.Get(service, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		return 1
	}

	// stdout, no trailing newline
	fmt.Print(val)
	return 0
}
