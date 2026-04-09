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

	if err != nil {
		return Location{}, fmt.Errorf("Error reading body: %w", err)
	}

	if res.StatusCode > 299 {
		return Location{}, fmt.Errorf("Response failed with status code: %d and \nBody: %s", res.StatusCode, body)
	}

	cache.Add(url, body)

	location := Location{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		return Location{}, fmt.Errorf("Problem unmarshalling location data: %w", err)
	}

	return location, nil
}

type AreaInfo struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func fetchPokemonInArea(area string, cache *pokecache.Cache) (AreaInfo, error) {
	const baseURL = "https://pokeapi.co/api/v2/location-area/"
	fullURL := baseURL + area + "/"

	var data []byte
	entry, ok := cache.Get(fullURL)
	if ok {
		data = entry
	} else {
		res, err := http.Get(fullURL)
		if err != nil {
			return AreaInfo{}, fmt.Errorf("Error calling API: %w", err)
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		if err != nil {
			return AreaInfo{}, fmt.Errorf("Error reading body: %w", err)
		}

		if res.StatusCode > 299 {
			return AreaInfo{}, fmt.Errorf("Response failed with status code: %d and \nBody: %s", res.StatusCode, body)
		}

		cache.Add(fullURL, body)
		data = body
	}
	areainfo := AreaInfo{}
	err := json.Unmarshal(data, &areainfo)
	if err != nil {
		return AreaInfo{}, err
	}

	return areainfo, nil
}

type PokemonInfo struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func fetchPokemonInfo(name string, cache *pokecache.Cache) (PokemonInfo, error) {
	const baseURL = "https://pokeapi.co/api/v2/pokemon/"
	fullURL := baseURL + name + "/"

	var data []byte
	entry, ok := cache.Get(fullURL)
	if ok {
		data = entry
	} else {
		res, err := http.Get(fullURL)
		if err != nil {
			return PokemonInfo{}, fmt.Errorf("Error calling API: %w", err)
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		if err != nil {
			return PokemonInfo{}, fmt.Errorf("Error reading body: %w", err)
		}

		if res.StatusCode > 299 {
			return PokemonInfo{}, fmt.Errorf("Response failed with status code: %d and \nBody: %s", res.StatusCode, body)
		}

		cache.Add(fullURL, body)
		data = body
	}
	pokemonStats := PokemonInfo{}
	err := json.Unmarshal(data, &pokemonStats)
	if err != nil {
		return PokemonInfo{}, err
	}

	return pokemonStats, nil
}
