// internal/pokeapi/pokeapi.go
// (boots) https://www.boot.dev/lessons/813eafe1-2e1d-42a0-b358-53e0f4d4fdc8
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/neixir/pokedex/internal/pokecache"
)

const LocationAreaUrl = "https://pokeapi.co/api/v2/location-area/"

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaInfo struct {
	EncounterMethodRates []EncounterMethodRates `json:"encounter_method_rates"`
	GameIndex            int                    `json:"game_index"`
	ID                   int                    `json:"id"`
	Location             Location               `json:"location"`
	Name                 string                 `json:"name"`
	Names                []Names                `json:"names"`
	PokemonEncounters    []PokemonEncounters    `json:"pokemon_encounters"`
}
type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

//	type VersionDetails struct {
//		Rate    int     `json:"rate"`
//		Version Version `json:"version"`
//	}
type EncounterMethodRates struct {
	EncounterMethod EncounterMethod `json:"encounter_method"`
	//VersionDetails  []VersionDetails `json:"version_details"`
	VersionDetails []struct {
		Rate    int     `json:"rate"`
		Version Version `json:"version"`
	} `json:"version_details"`
}
type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Names struct {
	Language Language `json:"language"`
	Name     string   `json:"name"`
}
type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Method struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type EncounterDetails struct {
	Chance          int    `json:"chance"`
	ConditionValues []any  `json:"condition_values"`
	MaxLevel        int    `json:"max_level"`
	Method          Method `json:"method"`
	MinLevel        int    `json:"min_level"`
}
type VersionDetails struct {
	EncounterDetails []EncounterDetails `json:"encounter_details"`
	MaxChance        int                `json:"max_chance"`
	Version          Version            `json:"version"`
}
type PokemonEncounters struct {
	Pokemon        Pokemon          `json:"pokemon"`
	VersionDetails []VersionDetails `json:"version_details"`
}

func GetLocationArea(url string, cache *pokecache.Cache) (LocationArea, error) {
	area := LocationArea{}

	// Si es al cache ho retornem
	body, ok := cache.Get(url)
	if ok {
		fmt.Printf("Obtenint %s del cache.\n", url)
	} else {
		// https://pkg.go.dev/net/http#example-Get
		res, err := http.Get(url)
		if err != nil {
			return area, fmt.Errorf("could not connect to PokeAPI")
		}

		body, err = io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			return area, fmt.Errorf("response failed with status code: %d", res.StatusCode)
		}
		if err != nil {
			return area, err
		}

		cache.Add(url, body)
		fmt.Printf("Afegint %s al cache.\n", url)

		defer res.Body.Close()
	}

	// https://blog.boot.dev/golang/json-golang/#example-unmarshal-json-to-struct-decode
	err := json.Unmarshal(body, &area)
	if err != nil {
		return area, err
	}

	return area, nil
}

func GetPokemonByArea(areaName string, cache *pokecache.Cache) ([]string, error) {
	names := []string{}

	url := fmt.Sprintf("%s%s", LocationAreaUrl, areaName)

	// Si es al cache ho retornem
	body, ok := cache.Get(areaName)
	if ok {
		fmt.Printf("Obtenint %s del cache.\n", areaName)
	} else {
		// https://pkg.go.dev/net/http#example-Get
		res, err := http.Get(url)
		if err != nil {
			return names, fmt.Errorf("could not connect to PokeAPI")
		}

		body, err = io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			if res.StatusCode == 404 {
				return names, fmt.Errorf("response failed with status code: %d (probably no area with that name)", res.StatusCode)
			}
			return names, fmt.Errorf("response failed with status code: %d", res.StatusCode)
		}
		if err != nil {
			return names, err
		}

		cache.Add(url, body)
		fmt.Printf("Afegint %s al cache.\n", areaName)

		defer res.Body.Close()
	}

	// https://blog.boot.dev/golang/json-golang/#example-unmarshal-json-to-struct-decode
	areaInfo := LocationAreaInfo{}
	err := json.Unmarshal(body, &areaInfo)
	if err != nil {
		return names, err
	}

	for i, _ := range areaInfo.PokemonEncounters {
		names = append(names, areaInfo.PokemonEncounters[i].Pokemon.Name)
	}

	return names, nil
}
