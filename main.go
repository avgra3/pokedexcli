package main

import (
	"bufio"
	"errors"
	"fmt"
	"maps"
	"os"
	"strings"
	"time"

	pokeapi "github.com/avgra3/pokedexcli/internal/pokeapi"
	pokecache "github.com/avgra3/pokedexcli/internal/pokecache"
	pokecatch "github.com/avgra3/pokedexcli/internal/pokecatch"
)

func main() {
	userPokedex := make(map[string]pokeapi.Pokemon)
	configuration := config{}
	configuration.UserPokedex = userPokedex
	interval := time.Second * 60
	cachePointer := pokecache.NewCache(interval)

	commands := map[string]cliCommand{
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
			description: "The 20 locations in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "The previous 20 locations in the Pokemon world",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore <LOCATION_NAME>",
			description: "See all Pokemon at a given location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <POKEMON_NAME>",
			description: "Attempt to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <POKEMON_NAME>",
			description: "Will return the name, heigjt, weight, stats and type(s) of the pokemon.",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "See all Pokemon currently in your pokedex",
			callback:    commandPokedex,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleaned := strings.Fields(strings.ToLower(input))
		value, ok := commands[cleaned[0]]
		if !ok {
			fmt.Println("Unknown command")
		}
		if ok {
			searchTerm := ""
			if len(cleaned) > 1 {
				searchTerm = cleaned[1]
			}
			err := value.callback(&configuration, cachePointer, searchTerm)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func cleanInput(text string) []string {
	cleanedText := strings.TrimSpace(text)
	cleanedText = strings.ToLower(cleanedText)
	return strings.Fields(cleanedText)
}

func commandExit(configuration *config, cache *pokecache.Cache, input string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(configuration *config, cache *pokecache.Cache, input string) error {
	message := fmt.Sprintf("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\nexplore <LOCATION_NAME>: Display all pokemon at a given location.\ncatch <POKEMON_NAME>: Attempt to catch a new pokemon. New Pokemon are added to the user's Pokedex\npokedex: See all Pokemon currently in your pokedex.")
	fmt.Println(message)
	return nil
}

func commandMap(configuration *config, cache *pokecache.Cache, input string) error {
	// Get 20 location areas in the Pokemon world
	// Each subsequent call gets the next 20 locations
	const POKEAPI = "https://pokeapi.co/api/v2/location-area"
	locationsResult, err := pokeapi.GetLocations(POKEAPI, cache, input)
	if err != nil {
		return err
	}

	if configuration.Next != "" {
		_, err := pokeapi.GetLocations(configuration.Next, cache, input)
		if err != nil {
			return err
		}
	}
	configuration.Next = locationsResult.Next
	configuration.Previous = locationsResult.Previous
	// Allow for first page to just return the current page
	if locationsResult.Previous == "" {
		configuration.Previous = POKEAPI
	}
	for _, value := range locationsResult.Results {
		fmt.Println(value.Name)
	}
	return nil
}

func commandMapBack(configuration *config, cache *pokecache.Cache, input string) error {
	previousApiUrl := configuration.Previous
	if previousApiUrl != "" {
		locationsResult, err := pokeapi.GetLocations(previousApiUrl, cache, input)
		if err != nil {
			return err
		}
		configuration.Next = locationsResult.Next
		configuration.Previous = locationsResult.Previous
		if locationsResult.Previous == "" {
			configuration.Previous = previousApiUrl
		}

		for _, value := range locationsResult.Results {
			fmt.Println(value.Name)
		}
		return nil

	}
	e := errors.New("There is no \"previous\" page of locations")
	return e
}

func commandCatch(configuration *config, cache *pokecache.Cache, input string) error {

	attemptMessage := fmt.Sprintf("Throwing a Pokeball at %v...", input)
	fmt.Println(attemptMessage)

	// Get pokemon info
	url := "https://pokeapi.co/api/v2/pokemon/" + input
	pokemonInfo, err := pokeapi.GetPokemon(url, cache, input)
	if err != nil {
		return err
	}
	// Get chances of success
	successMin := pokemonInfo.BaseExperience

	// Need a success and failure message
	success := fmt.Sprintf("%v was caught!", input)
	failure := fmt.Sprintf("%v escaped!", input)
	caught := pokecatch.SuccessfulCatch(successMin)
	if caught {
		(*configuration).UserPokedex[pokemonInfo.Name] = pokemonInfo
		fmt.Println(success)
	} else {
		fmt.Println(failure)
	}

	return nil
}

func commandExplore(configuration *config, cache *pokecache.Cache, input string) error {
	// Base lookup location
	baseUrl := "https://pokeapi.co/api/v2/location-area/" + input
	// Now try to get location
	locationAreaDetails, err := pokeapi.GetLocationAreas(baseUrl, cache, input)
	if err != nil {
		return err
	}
	// Our slice of pokemon encounters
	pokemonEncounters := locationAreaDetails.PokemonEncounters
	// Print out our pokemon names
	fmt.Printf("Exploring %v...\n", input)
	for _, pokemon := range pokemonEncounters {
		fmt.Printf("- %v\n", pokemon.Pokemon.Name)
	}

	return nil
}

func commandInspect(configuration *config, cache *pokecache.Cache, input string) error {
	// Have a message that tells user if the Pokemon they are looking for does not exist
	pokemon, ok := (*configuration).UserPokedex[input]
	if !ok {
		outputMessage := fmt.Sprintf("you have not caught that Pokemon")
		fmt.Println(outputMessage)
		return nil
	}
	pokemonName := pokemon.Name
	pokemonHeight := pokemon.Height
	pokemonWeight := pokemon.Weight
	pokemonStats := pokemon.Stats
	pokemonHP := "\t- hp: "
	pokemonAttack := "\t- attack: "
	pokemonDefense := "\t- defense: "
	pokemonSpecialAttack := "\t- special-attack: "
	pokemonSpecialDefense := "\t- special-defense: "
	pokemonSpeed := "\t- speed: "
	for _, stat := range pokemonStats {
		switch stat.Stat.Name {
		case "hp":
			pokemonHP = fmt.Sprintf("%v %v", pokemonHP, stat.BaseStat)
		case "attack":
			pokemonAttack = fmt.Sprintf("%v %v", pokemonAttack, stat.BaseStat)
		case "defense":
			pokemonDefense = fmt.Sprintf("%v %v", pokemonDefense, stat.BaseStat)

		case "special-attack":
			pokemonSpecialAttack = fmt.Sprintf("%v %v", pokemonSpecialAttack, stat.BaseStat)
		case "special-defense":
			pokemonSpecialDefense = fmt.Sprintf("%v %v", pokemonSpecialDefense, stat.BaseStat)
		case "speed":
			pokemonSpeed = fmt.Sprintf("%v %v", pokemonSpeed, stat.BaseStat)
		}
	}
	pokemonTypes := pokemon.Types
	typeNames := []string{}
	for _, pokemonType := range pokemonTypes {
		typeNames = append(typeNames, pokemonType.Type.Name)
	}

	// Output message(s)
	fmt.Printf("Name: %v\n", pokemonName)
	fmt.Printf("Height: %v\n", pokemonHeight)
	fmt.Printf("Weight: %v\n", pokemonWeight)
	fmt.Println("Stats:")
	fmt.Println(pokemonHP)
	fmt.Println(pokemonAttack)
	fmt.Println(pokemonDefense)
	fmt.Println(pokemonSpecialAttack)
	fmt.Println(pokemonSpecialDefense)
	fmt.Println(pokemonSpeed)
	fmt.Println("Types:")
	for _, pokemonType := range typeNames {
		typeOut := fmt.Sprintf("\t- %v", pokemonType)
		fmt.Println(typeOut)
	}

	return nil
}

func commandPokedex(configuration *config, cache *pokecache.Cache, input string) error {

	currentPokedex := (*configuration).UserPokedex
	pokemonNames := maps.Keys(currentPokedex)
	fmt.Println("Your Pokedex:")
	for key := range pokemonNames {
		fmt.Printf("\t- %v\n", key)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache, string) error
}

type config struct {
	Next        string
	Previous    string
	UserPokedex map[string]pokeapi.Pokemon
}
