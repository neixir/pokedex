package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/neixir/pokedex/internal/pokeapi"
	"github.com/neixir/pokedex/internal/pokecache"
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
	Next               *string
	Previous           *string
	Argv               []string
	locationNamesCache *pokecache.Cache
	pokemonNamesCache  *pokecache.Cache
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
	// Aixo de PokeApiUrl ho haura de fer a pokeapi.go
	url := PokeApiUrl
	if config.Previous != nil {
		url = *config.Next
	}

	area, err := pokeapi.GetLocationArea(url, config.locationNamesCache)
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
	// Aixo de PokeApiUrl ho haura de fer a pokeapi.go
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	url := config.Previous

	area, err := pokeapi.GetLocationArea(*url, config.locationNamesCache)
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

func commandExplore(config *Config) error {
	var areaName string

	if len(config.Argv) >= 2 {
		areaName = config.Argv[1]
	} else {
		return fmt.Errorf("missing parameter <area name>")
	}

	fmt.Printf("Exploring %s...\n", areaName)

	names, err := pokeapi.GetPokemonByArea(areaName, config.locationNamesCache)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, name := range names {
		fmt.Printf("- %s\n", name)
	}

	return nil

}

var supportedCommands = map[string]cliCommand{}

func main() {
	config := Config{
		locationNamesCache: pokecache.NewCache(5 * time.Second),
		pokemonNamesCache:  pokecache.NewCache(20 * time.Second),
	}

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
		"explore": {
			name:        "explore",
			description: "Lists all the pokemon located in an area",
			callback:    commandExplore,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := scanner.Text()
		config.Argv = cleanInput(input)

		if len(config.Argv) > 0 {
			command, ok := supportedCommands[config.Argv[0]]
			if ok {
				// clean[1] nomes funcionara amb "explore", amb la resta petara...
				// el que hauriem de fer es posar clean a config
				// i dins de la funcio agafar mes elements si cal
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
