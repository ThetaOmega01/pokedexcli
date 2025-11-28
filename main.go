package main

import (
	"bufio"
	"fmt"
	"os"

	"pokedexcli/cli"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Display prompt
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			continue
		}

		// Read user input and process command
		text := scanner.Text()
		callback, args, ok := cli.ProcessCommand(text)
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := callback(args...)
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}
