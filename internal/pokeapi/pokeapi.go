package pokeapi

import (
	"encoding/json"
	"fmt"
	pokecache "github.com/avgra3/pokedexcli/internal/pokecache"
	"io"
	"log"
	"net/http"
)

func GetLocationAreas(url string, cache *pokecache.Cache, name string) (LocationArea, error) {
	// Url+Name
	urlName := url
	// Check cache has data
	if cacheEntry, ok := cache.Get(urlName); ok {
		fmt.Println("We are using the cache...")
		// We get back a []bytes which we would want to convert to a location result
		cachedResult := LocationArea{}
		err := json.Unmarshal(cacheEntry, &cachedResult)
		if err != nil {
			log.Fatal(err)
		}
		return cachedResult, nil
	}

	// Cache has no data, call the api
	res, err := http.Get(urlName)
	if err != nil {
		return LocationArea{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	// Add new data to cache
	cache.Add(urlName, body)

	if res.StatusCode > 299 {
		statusCodeMessage := fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		log.Fatal(statusCodeMessage)
	}
	if err != nil {
		return LocationArea{}, err
	}
	locationResults := LocationArea{}
	err = json.Unmarshal(body, &locationResults)
	if err != nil {
		return LocationArea{}, err
	}
	return locationResults, nil
}

func GetLocations(url string, cache *pokecache.Cache, name string) (LocationResult, error) {
	// Check cache has data
	if cacheEntry, ok := cache.Get(url); ok {
		fmt.Println("We are using the cache...")
		// We get back a []bytes which we would want to convert to a location result
		cachedResult := LocationResult{}
		err := json.Unmarshal(cacheEntry, &cachedResult)
		if err != nil {
			return LocationResult{}, err
		}
		return cachedResult, nil

	}

	// Cache has no data, call the api
	res, err := http.Get(url)
	if err != nil {
		return LocationResult{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	// Add new data to cache
	cache.Add(url, body)

	if res.StatusCode > 299 {
		statusCodeMessage := fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return LocationResult{}, fmt.Errorf(statusCodeMessage)
	}
	if err != nil {
		log.Fatal(err)
	}
	locationResults := LocationResult{}
	err = json.Unmarshal(body, &locationResults)
	if err != nil {
		return LocationResult{}, err
	}
	return locationResults, nil
}

func GetPokemon(url string, cache *pokecache.Cache, name string) (Pokemon, error) {
	// Check cache has data
	if cacheEntry, ok := cache.Get(url); ok {
		fmt.Println("We are using the cache...")
		// We get back a []bytes which we would want to convert to a location result
		cachedResult := Pokemon{}
		err := json.Unmarshal(cacheEntry, &cachedResult)
		if err != nil {
			return Pokemon{}, err
		}
		return cachedResult, nil

	}

	// Cache has no data, call the api
	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	// Add new data to cache
	cache.Add(url, body)

	if res.StatusCode > 299 {
		statusCodeMessage := fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return Pokemon{}, fmt.Errorf(statusCodeMessage)
	}
	if err != nil {
		log.Fatal(err)
	}
	pokemonResult := Pokemon{}
	err = json.Unmarshal(body, &pokemonResult)
	if err != nil {
		return Pokemon{}, err
	}
	return pokemonResult, nil
}
