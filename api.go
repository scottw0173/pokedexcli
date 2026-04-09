package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

type Location struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func fetchLocations(url string, cache *pokecache.Cache) (Location, error) {
	data, ok := cache.Get(url)
	if ok {
		location := Location{}
		err := json.Unmarshal(data, &location)
		if err != nil {
			return Location{}, err
		}
		return location, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return Location{}, fmt.Errorf("Error calling API: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	cache.Add(url, body)

	if err != nil {
		return Location{}, fmt.Errorf("Error reading body: %w", err)
	}

	if res.StatusCode > 299 {
		return Location{}, fmt.Errorf("Response failed with status code: %d and \nBody: %s", res.StatusCode, body)
	}

	location := Location{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		return Location{}, fmt.Errorf("Problem unmarshalling location data: %w", err)
	}

	return location, nil
}
