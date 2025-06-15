package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

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

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))

	return words
}

func getLocationArea(url string) (LocationArea, error) {
	area := LocationArea{}

	// https://pkg.go.dev/net/http#example-Get
	res, err := http.Get(url)
	if err != nil {
		return area, fmt.Errorf("could not connect to PokeAPI")
	}

	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		return area, fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}
	if err != nil {
		return area, err
	}

	defer res.Body.Close()

	// https://blog.boot.dev/golang/json-golang/#example-unmarshal-json-to-struct-decode
	area = LocationArea{}
	err = json.Unmarshal(body, &area)
	if err != nil {
		return area, err
	}

	return area, nil
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
	url := "https://pokeapi.co/api/v2/location-area/"
	if config.Previous != nil {
		url = *config.Next
	}

	area, err := getLocationArea(url)
	if err != nil {
		return nil
	}

	// Actualitzem next i previous
	config.Previous = &url
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

	area, err := getLocationArea(*url)
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
