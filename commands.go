package main

import (
	"fmt"
	"math/rand"
	"os"
	"pokedexcli/internal/pokecache"
)

type config struct {
	next       string
	previous   string
	hasFetched bool
	pokedex    map[string]PokemonInfo
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache, []string) error
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
		"explore": {
			name:        "explore",
			description: "lists pokemon found in area typed after the word 'explore'",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "throws Pokeball in an attempt to catch named pokemon and add to user's Pokedex",
			callback:    commandCatch,
		},
	}
}

func commandExit(cfg *config, cache *pokecache.Cache, area []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, cache *pokecache.Cache, area []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config, cache *pokecache.Cache, area []string) error {
	var url string
	if !cfg.hasFetched {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else if cfg.next == "" {
		fmt.Println("No more locations.")
		return nil
	} else {
		url = cfg.next
	}

	locations, err := fetchLocations(url, cache)
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

func commandMapBack(cfg *config, cache *pokecache.Cache, area []string) error {
	if cfg.previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	locations, err := fetchLocations(cfg.previous, cache)
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

func commandExplore(cfg *config, cache *pokecache.Cache, area []string) error {
	if len(area) < 2 {
		return fmt.Errorf("please give name of an area")
	}

	areaName := area[1]
	areaInfo, _ := fetchPokemonInArea(areaName, cache)
	for _, pokemon := range areaInfo.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config, cache *pokecache.Cache, words []string) error {
	if len(words) < 2 {
		return fmt.Errorf("please name a Pokemon to try to catch")
	}

	pokemon := words[1]
	stats, _ := fetchPokemonInfo(pokemon, cache)

	catchDifficulty := stats.BaseExperience
	playerExperience := len(cfg.pokedex)
	randomness := rand.Intn(stats.BaseExperience)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
	successThreshold := randomness + ((playerExperience + 1) * 30)
	if successThreshold > catchDifficulty {
		cfg.pokedex[pokemon] = stats
		fmt.Printf("%s was caught!\n", pokemon)
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}
	return nil
}
