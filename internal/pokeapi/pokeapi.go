package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Zmahl/pokedexcli/internal/pokecache"
)

type Locations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokemonInLocation struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
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
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func (c *Client) GetLocations(pageURL *string, cache *pokecache.Cache) (Locations, error) {
	url := "https://pokeapi.co/api/v2/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	locations := Locations{}
	locationData, exists := cache.Get(url)
	if exists {
		err := json.Unmarshal(locationData, &locations)
		if err != nil {
			return Locations{}, err
		}

		return locations, nil
	}

	// Creates a request with the data provided
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Locations{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Locations{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Locations{}, err
	}

	err = json.Unmarshal(data, &locations)
	if err != nil {
		return Locations{}, err
	}

	cache.Add(url, data)

	return locations, nil
}

func (c *Client) Explore(area string, cache *pokecache.Cache) (PokemonInLocation, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", area)

	pokemonLocations := PokemonInLocation{}
	pokemonData, exists := cache.Get(area)
	if exists {
		err := json.Unmarshal(pokemonData, &pokemonLocations)
		if err != nil {
			return PokemonInLocation{}, err
		}
		return pokemonLocations, nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokemonInLocation{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonInLocation{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonInLocation{}, err
	}

	err = json.Unmarshal(data, &pokemonLocations)
	if err != nil {
		return PokemonInLocation{}, err
	}

	return pokemonLocations, nil
}

func (c *Client) GetPokemon(pokemon string) (PokemonInfo, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", pokemon)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokemonInfo{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonInfo{}, err
	}
	defer resp.Body.Close()

	pokemonInfo := PokemonInfo{}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonInfo{}, err
	}

	err = json.Unmarshal(data, &pokemonInfo)
	if err != nil {
		return PokemonInfo{}, err
	}

	return pokemonInfo, nil
}
