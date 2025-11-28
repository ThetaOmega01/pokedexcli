package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"pokedexcli/cache"
)

var (
	nextURL     = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	previousURL = ""
	fetchCache  = cache.NewCache(1 * time.Minute)
)

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ResourceList struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Result `json:"results"`
}

func fetchResourceList(url string) (*ResourceList, error) {
	// Check cache first
	if cachedData, ok := fetchCache.Get(url); ok {
		var resourceList ResourceList
		if err := json.Unmarshal(cachedData, &resourceList); err != nil {
			return nil, err
		}

		previousURL = resourceList.Previous
		nextURL = resourceList.Next

		return &resourceList, nil
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

	// Read and cache the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Add to cache
	fetchCache.Add(url, body)

	// Decode the response body
	var resourceList ResourceList
	if err := json.Unmarshal(body, &resourceList); err != nil {
		return nil, err
	}

	previousURL = resourceList.Previous
	nextURL = resourceList.Next

	return &resourceList, nil
}

func printLocationAreas(resourceList *ResourceList) {
	for _, locationArea := range resourceList.Results {
		fmt.Println(locationArea.Name)
	}
}

func commandMap(args ...string) error {
	resourceList, err := fetchResourceList(nextURL)
	if err != nil {
		return err
	}

	printLocationAreas(resourceList)

	return nil
}

func commandMapBack(args ...string) error {
	// Handle case where there is no previous URL
	if previousURL == "" {
		fmt.Println("No previous location areas available.")
		return nil
	}

	resourceList, err := fetchResourceList(previousURL)
	if err != nil {
		return err
	}

	printLocationAreas(resourceList)

	return nil
}
