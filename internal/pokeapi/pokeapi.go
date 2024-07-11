package pokeapi

import (
	"encoding/json"
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
