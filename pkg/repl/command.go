package repl

import (
	"errors"
	"fmt"
	"os"

	"github.com/Zmahl/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(string, *Config, *pokecache.Cache) error
}

func CheckCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "View the next set of location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "View the previous set of location areas",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "View all pokemon in a given area",
			callback:    commandExplore,
		},
	}
}

func commandHelp(area string, c *Config, cache *pokecache.Cache) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, c := range CheckCommands() {
		fmt.Printf("%s: %s", c.name, c.description)
		fmt.Println()
	}

	fmt.Println()
	fmt.Println()

	return nil
}

func commandExit(area string, config *Config, cache *pokecache.Cache) error {
	os.Exit(0)
	return nil
}

func commandMap(area string, config *Config, cache *pokecache.Cache) error {
	locations, err := config.PokeApiClient.GetLocations(config.Next, cache)
	if err != nil {
		return err
	}

	config.Next = locations.Next
	config.Previous = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapB(area string, config *Config, cache *pokecache.Cache) error {
	if config.Previous == nil {
		return errors.New("this is the first page")
	}

	locations, err := config.PokeApiClient.GetLocations(config.Previous, cache)
	if err != nil {
		return err
	}

	config.Next = locations.Next
	config.Previous = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandExplore(area string, config *Config, cache *pokecache.Cache) error {
	pokemonLocations, err := config.PokeApiClient.Explore(area, cache)
	if err != nil {
		errMessage := fmt.Sprintf("could not find area: %v", area)
		return errors.New(errMessage)
	}
	fmt.Printf("Exploring area %v...\n", area)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range pokemonLocations.PokemonEncounters {
		fmt.Printf("  - %v\n", pokemon.Pokemon.Name)
	}

	return nil
}
