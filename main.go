package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
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
	// I used a map[string]Pokemon to keep track of caught Pokemon.
	caughtPokemon map[string]pokeapi.PokemonType
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

	names, err := pokeapi.GetPokemonNamesByArea(areaName, config.locationNamesCache)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, name := range names {
		fmt.Printf("- %s\n", name)
	}

	return nil

}

func commandCatch(config *Config) error {
	var pokemonName string

	if len(config.Argv) >= 2 {
		pokemonName = config.Argv[1]
	} else {
		return fmt.Errorf("missing parameter <pokemon name>")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := pokeapi.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	// fmt.Printf("Trying to catch %s (base experience %d).\n", pokemon.Name, pokemon.BaseExperience)

	// You can use the pokemon's "base experience" to determine the chance of catching it.
	// The higher the base experience, the harder it should be to catch.
	// https://claude.ai/chat/b741ba22-fbfa-4a87-9ef5-02335c9a5bfd
	// Power Decay
	probability := int(100 / math.Pow(float64(pokemon.BaseExperience), 0.2))
	random := rand.Intn(100)
	if random < probability {
		fmt.Printf("%s was caught!\n", pokemonName)
		// fmt.Printf("%v < %v\n", random, probability)
		// Once the Pokemon is caught, add it to the user's Pokedex.
		config.caughtPokemon[pokemonName] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
		// fmt.Printf("%v >= %v\n", random, probability)
	}

	return nil

}

func commandInspect(config *Config) error {
	var pokemonName string

	if len(config.Argv) >= 2 {
		pokemonName = config.Argv[1]
	} else {
		return fmt.Errorf("missing parameter <pokemon name>")
	}

	pokemon, ok := config.caughtPokemon[pokemonName]
	if ok {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, typ := range pokemon.Types {
			fmt.Printf("  -%v\n", typ.Type.Name)
		}
	} else {
		fmt.Println("you have not caught that pokemon")
	}

	return nil
}

// Mostrem els pokemons que s'han obtingut
func commandPokedex(config *Config) error {
	if len(config.caughtPokemon) > 0 {
		fmt.Println("Your Pokedex:")
		for _, pokemon := range config.caughtPokemon {
			fmt.Printf("- %s (%d XP)\n", pokemon.Name, pokemon.BaseExperience)
		}
	} else {
		fmt.Println("Your Pokedex is empty :(")
	}

	return nil
}

var supportedCommands = map[string]cliCommand{}

func main() {
	config := Config{
		locationNamesCache: pokecache.NewCache(5 * time.Second),
		pokemonNamesCache:  pokecache.NewCache(20 * time.Second),
		caughtPokemon:      map[string]pokeapi.PokemonType{},
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

		// C2 L1 https://www.boot.dev/lessons/813eafe1-2e1d-42a0-b358-53e0f4d4fdc8
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},

		// C2 L1 https://www.boot.dev/lessons/813eafe1-2e1d-42a0-b358-53e0f4d4fdc8
		"mapb": {
			name:        "mapb",
			description: "Displays the 20 names of 20 previous location areas in the Pokemon world",
			callback:    commandMapB,
		},

		// C2 L3 https://www.boot.dev/lessons/e53abbb4-5d8a-4feb-ba08-828f03311e51
		"explore": {
			name:        "explore",
			description: "Lists all the pokemon located in an area",
			callback:    commandExplore,
		},

		// C2 L4 https://www.boot.dev/lessons/ed962683-cb2d-4989-99e9-5cfa144810b5
		"catch": {
			name:        "catch",
			description: "Catching Pokemon adds them to the user's Pokedex",
			callback:    commandCatch,
		},

		// C2 L5 https://www.boot.dev/lessons/0911b406-0b43-4bfe-b60c-177d859093e1
		"inspect": {
			name:        "inspect",
			description: "Prints the name, height, weight, stats and type(s) of the Pokemon",
			callback:    commandInspect,
		},

		// C3 L1 https://www.boot.dev/lessons/104a68ca-cea7-42ef-9321-fb8270000db2
		"pokedex": {
			name:        "pokedex",
			description: "Prints a list of all the names of the Pokemon the user has caught",
			callback:    commandPokedex,
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
