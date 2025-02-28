package pokeapi

import (
	"encoding/json"
	"fmt"
	pokecache "github.com/avgra3/pokedexcli/internal/pokecache"
	"io"
	"log"
	"net/http"
)

func GetLocations(url string, cache *pokecache.Cache) LocationResult {
	// Check cache has data
	if cacheEntry, ok := cache.Get(url); ok {
		fmt.Println("We are using the cache...")
		// We get back a []bytes which we would want to convert to a location result
		cachedResult := LocationResult{}
		err := json.Unmarshal(cacheEntry, &cachedResult)
		if err != nil {
			log.Fatal(err)
		}
		return cachedResult

	}

	// Cache has no data, call the api
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	// Add new data to cache
	cache.Add(url, body)

	if res.StatusCode > 299 {
		statusCodeMessage := fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		log.Fatal(statusCodeMessage)
	}
	if err != nil {
		log.Fatal(err)
	}
	locationResults := LocationResult{}
	err = json.Unmarshal(body, &locationResults)
	if err != nil {
		log.Fatal(err)
	}
	return locationResults
}

type Locations struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationResult struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []Locations `json:"results"`
}

type Generation struct {
}

type Pokedex struct{}

type VersionGroup struct{}

type Region struct {
	Id             int            `json:"id"`
	Locations      []Locations    `json:"locations"`
	Name           string         `json:"name"`
	Names          []LanguageName `json:"names"`
	MainGeneration Generation     `json:"main_generation"`
	Pokedexes      []Pokedex      `json:"pokedexes"`
	VersionGroups  []VersionGroup `json:"version_groups"`
}

type GenerationGameIndex struct{}

type Location struct {
	Id           int                   `json:"id"`
	Name         string                `json:"name"`
	Region       Region                `json:"region"`
	Names        []LanguageName        `json:"names"`
	GameIndicies []GenerationGameIndex `json:"game_indicies"`
	Areas        []LocationArea        `json:"areas"`
}

type LocationArea struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []EncounterMethodRates
	Location             Location
}

type SpecificPokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonEncounters struct {
	EncounterMethodRates []EncounterMethodRates `json:"encounter_method_rates"`
	GameIndex            int                    `json:"game_index"`
	Id                   int                    `json:"id"`
	Location             Locations              `json:"location"`
	LocationName         string                 `json:"name"`
	Names                []LanguageName         `json:"names"`
	PokemonEncounters    []PokemonEncounter     `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon        Pokemon          `json:"pokemon"`
	VersionDetails []VersionDetails `json:"version_details"`
}

type Pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LanguageName struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

type EncounterMethodRates struct {
	EncounterMethod         EncounterMethod          `json:"encounter_method"`
	VersionEncounterDetails []VersionEncounterDetail `json:"version_details"`
}

type VersionEncounterDetail struct {
	Version          Version            `json:"version"`
	MaxChance        int                `json:"chance"`
	EncounterDetails []EncounterDetails `json:"encoutner_details"`
}

type EncounterDetails struct {
	MinLevel        int                     `json:"min_level"`
	MaxLevel        int                     `json:"max_level"`
	Conditionvalues EncounterConditionValue `json:"condition_values"`
	Chance          int                     `json:"chance"`
	Method          EncounterMethod         `json:"method"`
}

type EncounterMethod struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type EncounterCondition struct {
	Id     int                       `json:"id"`
	Name   string                    `json:"name"`
	Names  []LanguageName            `json:"names"`
	Values []EncounterConditionValue `json:"values"`
}

type EncounterConditionValue struct {
	Id        int                `json:"id"`
	Name      string             `json:"name"`
	Condition EncounterCondition `json:"condition"`
	Names     []LanguageName     `json:"names"`
}

type VersionDetails struct {
	Rate    int     `json:"rate"`
	Version Version `json:"version"`
}

type Version struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
