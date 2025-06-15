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

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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
