package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func fetchLocations(url string) (Location, error) {
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

	location := Location{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		return Location{}, fmt.Errorf("Problem unmarshalling location data: %w", err)
	}

	return location, nil
}
