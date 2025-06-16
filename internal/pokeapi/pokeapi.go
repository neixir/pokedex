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
const PokemonUrl = "https://pokeapi.co/api/v2/pokemon/"

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

func GetPokemonNamesByArea(areaName string, cache *pokecache.Cache) ([]string, error) {
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

// TODO Utilitzar cache
func GetPokemon(name string) (PokemonType, error) {
	pokemon := PokemonType{}
	url := PokemonUrl + name

	res, err := http.Get(url)
	if err != nil {
		return pokemon, fmt.Errorf("could not connect to PokeAPI")
	}

	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		if res.StatusCode == 404 {
			return pokemon, fmt.Errorf("response failed with status code: %d (probably no pokemon with that name)", res.StatusCode)
		}
		return pokemon, fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}
	if err != nil {
		return pokemon, err
	}

	defer res.Body.Close()

	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return pokemon, err
	}

	return pokemon, nil
}
