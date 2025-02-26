package main

import (
	"bufio"
	"errors"
	"fmt"
	internal "github.com/avgra3/pokedexcli/internal"
	"os"
	"strings"
)

func main() {
	configuration := config{}

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
			err := value.callback(&configuration)
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

func commandExit(configuration *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(configuration *config) error {
	message := fmt.Sprintf("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\n")
	fmt.Print(message)
	return nil
}

func commandMap(configuration *config) error {
	// Get 20 location areas in the Pokemon world
	// Each subsequent call gets the next 20 locations
	locationsResult := internal.GetLocations("https://pokeapi.co/api/v2/location-area")

	if configuration.Next != "" {
		locationsResult = internal.GetLocations(configuration.Next)
	}
	configuration.Next = locationsResult.Next
	configuration.Previous = locationsResult.Previous
	for _, value := range locationsResult.Results {
		fmt.Println(value.Name)
	}
	return nil
}

func commandMapBack(configuration *config) error {
	previousApiUrl := configuration.Previous
	if previousApiUrl != "" {
		locationsResult := internal.GetLocations(previousApiUrl)
		configuration.Next = locationsResult.Next
		configuration.Previous = locationsResult.Previous

		for _, value := range locationsResult.Results {
			fmt.Println(value.Name)
		}
		return nil

	}
	e := errors.New("There is no \"previous\" page of locations")
	return e
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next     string
	Previous string
}
