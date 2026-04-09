package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokecache"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		pokedex: make(map[string]PokemonInfo),
	}
	cache := pokecache.NewCache(5 * time.Second)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		text := scanner.Text()

		words := cleanInput(text)
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		command, exists := commands[commandName]
		if exists {
			err := command.callback(cfg, cache, words)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
