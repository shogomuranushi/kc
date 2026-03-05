package cmd

import (
	"fmt"
	"os"
)

// subcommands that are handled internally
var subcommands = map[string]bool{
	"set":     true,
	"get":     true,
	"delete":  true,
	"list":    true,
	"run":     true,
	"migrate": true,
}

// Execute is the main entry point for the CLI.
func Execute() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		os.Exit(0)
	}

	command := args[0]
	cmdArgs := args[1:]

	var exitCode int

	if subcommands[command] {
		switch command {
		case "set":
			exitCode = RunSet(cmdArgs)
		case "get":
			exitCode = RunGet(cmdArgs)
		case "delete":
			exitCode = RunDelete(cmdArgs)
		case "list":
			exitCode = RunList(cmdArgs)
		case "run":
			exitCode = RunRun(cmdArgs)
		case "migrate":
			exitCode = RunMigrate(cmdArgs)
		}
	} else {
		// External command: run with .env expansion
		exitCode = RunExternalCommand(args)
	}

	os.Exit(exitCode)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `kc - Keychain CLI for secret management

Usage:
  kc <command> [args...]        Run command with .env secrets expanded
  kc set <service> <key> [val]  Save a secret to keychain
  kc get <service> <key>        Get a secret from keychain
  kc delete <service> <key>     Delete a secret from keychain
  kc list [service]             List stored secrets
  kc run [--env-file <path>] -- <command> [args...]
                                Run command with explicit options
  kc migrate [file]             Migrate plaintext .env to kc:// references`)
}
