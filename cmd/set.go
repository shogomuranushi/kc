package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/shogomuranushi/kc/internal/keychain"
	"github.com/shogomuranushi/kc/internal/uri"
	"golang.org/x/sys/unix"
)

var validNameChars = regexp.MustCompile(`^[a-zA-Z0-9_.\-]+$`)

func RunSet(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "[ERROR] Usage: kc set <service> <key> [value]")
		return 1
	}

	service := args[0]
	key := args[1]

	if !validNameChars.MatchString(service) || !validNameChars.MatchString(key) {
		fmt.Fprintln(os.Stderr, "[ERROR] Invalid characters in service/key (allowed: [a-zA-Z0-9_.-])")
		return 1
	}

	var value string
	if len(args) >= 3 {
		value = args[2]
	} else {
		// Check if stdin is a pipe
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Reading from pipe
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				value = strings.TrimSpace(scanner.Text())
			}
		} else {
			// Interactive prompt (no echo)
			fmt.Fprint(os.Stderr, "Enter value: ")
			b, err := readPassword()
			if err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] Failed to read input: %s\n", err)
				return 1
			}
			value = string(b)
			fmt.Fprintln(os.Stderr)
		}
	}

	if value == "" {
		fmt.Fprintln(os.Stderr, "[ERROR] Value cannot be empty")
		return 1
	}

	if err := keychain.Set(service, key, value); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		return 1
	}

	fmt.Fprintf(os.Stderr, "Saved: %s\n", uri.Format(service, key))
	return 0
}

func readPassword() ([]byte, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}

	newState := *oldState
	newState.Lflag &^= unix.ECHO
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &newState); err != nil {
		return nil, err
	}
	defer unix.IoctlSetTermios(fd, unix.TCSETS, oldState)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return []byte(strings.TrimRight(line, "\r\n")), nil
}
