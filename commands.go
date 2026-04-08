package main

import (
	"fmt"
	"os"
)

type config struct {
	next       string
	previous   string
	hasFetched bool
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "prints the first or next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "prints the previous 20 locations",
			callback:    commandMapBack,
		},
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	var url string
	if !cfg.hasFetched {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else if cfg.next == "" {
		fmt.Println("No more locations.")
		return nil
	} else {
		url = cfg.next
	}

	locations, err := fetchLocations(url)
	if err != nil {
		return err
	}

	cfg.next = locations.Next
	if locations.Previous != nil {
		cfg.previous = *locations.Previous
	} else {
		cfg.previous = ""
	}
	cfg.hasFetched = true

	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapBack(cfg *config) error {
	if cfg.previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	locations, err := fetchLocations(cfg.previous)
	if err != nil {
		return err
	}

	cfg.next = locations.Next
	if locations.Previous != nil {
		cfg.previous = *locations.Previous
	} else {
		cfg.previous = ""
	}
	cfg.hasFetched = true

	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	return nil
}
