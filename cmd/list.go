package cmd

import (
	"fmt"
	"os"

	"github.com/shogomuranushi/kc/internal/keychain"
)

func RunList(args []string) int {
	var service string
	if len(args) > 0 {
		service = args[0]
	}

	entries, err := keychain.List(service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		return 1
	}

	if len(entries) == 0 {
		fmt.Fprintln(os.Stderr, "No secrets found")
		return 0
	}

	for _, e := range entries {
		fmt.Fprintln(os.Stderr, e)
	}
	return 0
}
