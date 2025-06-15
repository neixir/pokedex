package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))

	return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, com := range supportedCommands {
		fmt.Printf("%s: %s\n", com.name, com.description)
	}
	return nil
}

var supportedCommands = map[string]cliCommand{}

func main() {

	supportedCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message", // "Lists all available commands",
			callback:    commandHelp,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		clean := cleanInput(input)
		// Amb aquest if mirem que no sigui linia en blanc,
		// pero no es suficient pel test
		// printf "help\nexit\n" | go run .
		if len(clean) > 0 {
			command, ok := supportedCommands[clean[0]]
			if ok {
				command.callback()
			} else {
				fmt.Println("Unknown command")
			}
		}

	}
}
