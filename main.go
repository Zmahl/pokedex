package main

import (
	"time"

	"github.com/Zmahl/pokedexcli/internal/pokeapi"
	"github.com/Zmahl/pokedexcli/internal/pokecache"
	"github.com/Zmahl/pokedexcli/pkg/repl"
)

func main() {
	pokeClient := pokeapi.NewClient(5 * time.Second)
	pokeCache := pokecache.NewCache(5 * time.Minute)
	config := &repl.Config{
		PokeApiClient: pokeClient,
		Pokedex:       make(map[string]pokeapi.PokemonInfo),
	}

	repl.StartRepl(config, pokeCache)
}
