// Package cli contains utils for command line processing and command types
package cli

import (
	"fmt"
	"os"
	"strings"
)

// Registry of supported commands
var supportedCommands map[string]cliCommand

func init() {
	supportedCommands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Show available commands",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

// start of core registered callbacks
func commandHelp(args ...string) error {
	helpText := "Welcome to the Pokedex!\nUseage:\n\n"

	for _, cmd := range supportedCommands {
		helpText += fmt.Sprintf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Print(helpText)

	return nil
}

func commandExit(args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// end of core registered callbacks

type cliCommand struct {
	name        string
	description string
	callback    func(args ...string) error
}

func ProcessCommand(input string) (func(...string) error, []string, bool) {
	// Split input into command and arguments
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, nil, false
	}

	cmdName := parts[0]
	args := parts[1:]

	command, ok := supportedCommands[cmdName]
	if !ok {
		return nil, nil, false
	}

	return command.callback, args, true
}
