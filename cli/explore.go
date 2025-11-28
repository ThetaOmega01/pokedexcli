package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
}

type LocationArea struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

func printPokemon(locationArea *LocationArea) {
	if len(locationArea.PokemonEncounters) == 0 {
		fmt.Println("No Pokemon found in this area.")
		return
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
}

func commandExplore(args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a location area name (usage: explore <area_name>)")
	}

	areaName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", areaName)

	fmt.Printf("Exploring %s...\n", areaName)
	// Check cache first
	if cachedData, ok := fetchCache.Get(url); ok {
		var locationArea LocationArea
		if err := json.Unmarshal(cachedData, &locationArea); err != nil {
			return err
		}

		printPokemon(&locationArea)
		return nil
	}

	// Data not in cache, fetch from API
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()

	if res.StatusCode == 404 {
		return fmt.Errorf("location area '%s' not found", areaName)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("error fetching location area: status %d", res.StatusCode)
	}

	// Read and cache the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Add to cache
	fetchCache.Add(url, body)

	// Decode the response body
	var locationArea LocationArea
	if err := json.Unmarshal(body, &locationArea); err != nil {
		return err
	}

	printPokemon(&locationArea)
	return nil
}
