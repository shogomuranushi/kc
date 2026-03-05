package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/shogomuranushi/kc/internal/dotenv"
	"github.com/shogomuranushi/kc/internal/keychain"
	"github.com/shogomuranushi/kc/internal/uri"
)

func RunRun(args []string) int {
	envFile := ""
	cmdArgs := args

	// Parse --env-file option
	for i := 0; i < len(cmdArgs); i++ {
		if cmdArgs[i] == "--env-file" && i+1 < len(cmdArgs) {
			envFile = cmdArgs[i+1]
			cmdArgs = append(cmdArgs[:i], cmdArgs[i+2:]...)
			break
		}
	}

	// Strip leading "--"
	if len(cmdArgs) > 0 && cmdArgs[0] == "--" {
		cmdArgs = cmdArgs[1:]
	}

	if len(cmdArgs) == 0 {
		fmt.Fprintln(os.Stderr, "[ERROR] Usage: kc run [--env-file <path>] -- <command> [args...]")
		return 1
	}

	return execWithEnv(envFile, cmdArgs)
}

// RunExternalCommand runs an external command with .env expansion.
func RunExternalCommand(args []string) int {
	if len(args) == 0 {
		return 1
	}
	return execWithEnv("", args)
}

func execWithEnv(envFilePath string, cmdArgs []string) int {
	// Resolve env file
	if envFilePath == "" {
		cwd, err := os.Getwd()
		if err == nil {
			envFilePath = dotenv.FindEnvFile(cwd)
		}
	}

	// Build environment
	env := os.Environ()

	if envFilePath != "" {
		entries, err := dotenv.Parse(envFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
			return 1
		}

		for _, e := range entries {
			if e.IsBlank || !e.HasValue {
				continue
			}

			val := e.Value
			if e.IsKCRef {
				parsed, err := uri.Parse(val)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
					return 1
				}
				secret, err := keychain.Get(parsed.Service, parsed.Key)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
					return 1
				}
				val = secret
			}

			env = append(env, fmt.Sprintf("%s=%s", e.Key, val))
		}
	}

	// Find command
	cmdPath, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Command not found: %s\n", cmdArgs[0])
		return 1
	}

	// Exec (replaces current process)
	err = syscall.Exec(cmdPath, cmdArgs, env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		return 1
	}
	return 0
}
