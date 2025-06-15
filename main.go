package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/neixir/pokedex/internal/pokeapi"
)

const PokeApiUrl = "https://pokeapi.co/api/v2/location-area/"

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

// This struct will contain the Next and Previous URLs that you'll need to paginate through location areas.
// CH2 L1 https://www.boot.dev/lessons/813eafe1-2e1d-42a0-b358-53e0f4d4fdc8
type Config struct {
	Next     *string
	Previous *string
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))

	return words
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, com := range supportedCommands {
		fmt.Printf("%s: %s\n", com.name, com.description)
	}
	return nil
}

func commandMap(config *Config) error {
	url := PokeApiUrl
	if config.Previous != nil {
		url = *config.Next
	}

	area, err := pokeapi.GetLocationArea(url)
	if err != nil {
		return nil
	}

	// Actualitzem next i previous
	config.Previous = &area.Previous
	config.Next = &area.Next

	// Mostrem els noms
	for _, loc := range area.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapB(config *Config) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	url := config.Previous

	area, err := pokeapi.GetLocationArea(*url)
	if err != nil {
		return nil
	}

	// Actualitzem next i previous
	config.Previous = &area.Previous
	config.Next = &area.Next

	//
	for _, loc := range area.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

var supportedCommands = map[string]cliCommand{}

func main() {
	config := Config{}

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
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the 20 names of 20 previous location areas in the Pokemon world",
			callback:    commandMapB,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := scanner.Text()
		clean := cleanInput(input)

		if len(clean) > 0 {
			command, ok := supportedCommands[clean[0]]
			if ok {
				err := command.callback(&config)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}
