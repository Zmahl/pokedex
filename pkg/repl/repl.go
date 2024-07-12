package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Zmahl/pokedexcli/internal/pokeapi"
	"github.com/Zmahl/pokedexcli/internal/pokecache"
)

type Config struct {
	PokeApiClient pokeapi.Client
	Pokedex       map[string]pokeapi.PokemonInfo
	Next          *string
	Previous      *string
}

func StartRepl(config *Config, cache *pokecache.Cache) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		command, exists := CheckCommands()[commandName]
		if exists {
			if len(words) >= 2 {
				err := command.callback(words[1], config, cache)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				err := command.callback("", config, cache)
				if err != nil {
					fmt.Println(err)
				}
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
