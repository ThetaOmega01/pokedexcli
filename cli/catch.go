package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
)

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

// User's Pokedex - stores caught Pokemon
var pokedex = make(map[string]Pokemon)

// calculateCatchChance returns a catch probability between 0 and 1
// Higher base experience makes it harder to catch
func calculateCatchChance(baseExperience int) float64 {
	// Base catch rate starts at 90% for low experience Pokemon
	// Decreases as base experience increases
	// Minimum catch rate is 10%

	if baseExperience <= 0 {
		return 0.9
	}

	// Formula: catch_rate = max(0.1, 0.9 - (base_exp / 500))
	// This means:
	// - 0 base_exp = 90% catch rate
	// - 100 base_exp = 70% catch rate
	// - 200 base_exp = 50% catch rate
	// - 400+ base_exp = 10% catch rate (minimum)

	catchRate := 0.9 - (float64(baseExperience) / 500.0)
	if catchRate < 0.1 {
		catchRate = 0.1
	}

	return catchRate
}

// attemptCatch returns true if the catch was successful
func attemptCatch(catchChance float64) bool {
	roll := rand.Float64()
	return roll < catchChance
}

func fetchPokemon(url string) (*Pokemon, error) {
	// Check cache first
	if cachedData, ok := fetchCache.Get(url); ok {
		var pokemon Pokemon
		if err := json.Unmarshal(cachedData, &pokemon); err != nil {
			return nil, err
		}
		return &pokemon, nil
	}

	// Data not in cache, fetch from API
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("pokemon not found")
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching pokemon: status %d", res.StatusCode)
	}

	// Read and cache the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Add to cache
	fetchCache.Add(url, body)

	// Decode the response body
	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return nil, err
	}

	return &pokemon, nil
}

func commandCatch(args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a Pokemon name (usage: catch <pokemon_name>)")
	}

	pokemonName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)

	// Fetch Pokemon data
	pokemon, err := fetchPokemon(url)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	// Calculate catch chance based on base experience
	// Higher base experience = harder to catch
	catchChance := calculateCatchChance(pokemon.BaseExperience)

	// Attempt to catch
	if attemptCatch(catchChance) {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		pokedex[pokemon.Name] = *pokemon
		fmt.Printf("You can now use inspect %s to view its details.\n", pokemon.Name)
		fmt.Printf("You now have %d Pokemon in your Pokedex.\n", len(pokedex))
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}
