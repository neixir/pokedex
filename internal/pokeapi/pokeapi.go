// internal/pokeapi/pokeapi.go
// (boots) https://www.boot.dev/lessons/813eafe1-2e1d-42a0-b358-53e0f4d4fdc8
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func GetLocationArea(url string) (LocationArea, error) {
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
