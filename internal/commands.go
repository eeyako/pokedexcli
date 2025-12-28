package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(string) error
}

/*
List of all commands availabe in the program.
Contains their name, description and Callback function
*/
func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays this message",
			Callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Lists the 20 areas forward available to 'explore'",
			Callback:    commandMapForward,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists the 20 areas backward available to 'explore'",
			Callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Lists all of the Pokemon of a given area",
			Callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Tries to catch a Pokemon!",
			Callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a given pokemon if present in the pokedex",
			Callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Prints out the Pokemons you have caught so far",
			Callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
	}
}

/* Callback functions below */

/*
commandHelp
Lists all comands available in the program for ease of use.
*/
func commandHelp(_ string) error {
	fmt.Println(`
Welcome to the Pokedex!
Usage:`)
	fmt.Println()
	for _, v := range GetCommands() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

/*
commandExit
Exits the program.
*/
func commandExit(_ string) error {
	os.Exit(0)
	return nil
}

/*
commandMapForward
Lists to the user the 20 next location areas to explore.
*/
func commandMapForward(_ string) error {
	if nextUrl == "" {
		fmt.Printf("Error: cannot map futher.\n")
		return nil
	}
	printLocationAreas(nextUrl)

	return nil
}

/*
commandMapBack
Lists to the user the 20 previous location areas to explore.
*/
func commandMapBack(_ string) error {
	if previousUrl == "" {
		fmt.Printf("Error: cannot map back.\n")
		return nil
	}
	printLocationAreas(previousUrl)

	return nil
}

/*
commandExplore
Explores the given area, listing all pokemons that are available for catching
*/
func commandExplore(areaName string) error {
	if areaName == "" {
		fmt.Printf("Please provide the name of the area to explore.\n")
		return nil
	}
	_, valid := currentLocationAreas[areaName]
	if !valid {
		fmt.Printf("\nCould not find %v in the current locations\n", areaName)
		return nil
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", areaName)

	// Fetch data from either cache or GET
	data, err := FetchData(url)
	if err != nil {
		return err
	}

	var locationarea LocationArea
	err = json.Unmarshal(*data, &locationarea)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nExploring %s...\n", locationarea.Name)
	fmt.Println("Found Pokemon:")
	// Populate only with pokemons that haven't been removed
	currentPokemons = make(map[string]bool)
	for i := range locationarea.PokemonEncounters {
		pokemonEncounter := locationarea.PokemonEncounters[i]
		_, ok := removedPokemons[areaName][pokemonEncounter.Pokemon.Name]
		if ok {
			continue
		}

		currentPokemons[pokemonEncounter.Pokemon.Name] = true
		fmt.Printf(" - %s\n", pokemonEncounter.Pokemon.Name)
	}
	currentLocationArea = areaName

	return nil
}

/*
commandCatch
Tries to catch the given Pokemon name. If succeedes, adds that to the pokemon
map
*/
func commandCatch(pokemonName string) error {
	if pokemonName == "" {
		fmt.Printf("Please provide the name of the Pokemon to try to catch.\n")
		return nil
	}
	if currentLocationArea == "" {
		fmt.Printf("You need to explore an area fist!\n")
		return nil
	}
	_, ok := pokedex[pokemonName]
	if ok {
		fmt.Printf("%v is already on the Pokedex!\n", pokemonName)
		return nil
	}
	_, ok = currentPokemons[pokemonName]
	if !ok {
		fmt.Printf("Could not find %v in %v\n", pokemonName, currentLocationArea)
		return nil
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", pokemonName)

	// Fetch data form either cache or GET
	data, err := FetchData(url)
	if err != nil {
		return err
	}

	var pokemon Pokemon
	err = json.Unmarshal(*data, &pokemon)
	if err != nil {
		log.Fatal(err)
	}

	// Try to catch Pokemon
	fmt.Printf("\nThrowing a Pokeball at %v", pokemonName)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Duration(1) * time.Second)
		fmt.Printf(".")
	}
	fmt.Printf(" ")
	success, ranAway := catchAttempt(pokemon)
	remove := false
	if success {
		fmt.Printf("it was %v!\n", color.GreenString("CAUGHT"))
		fmt.Printf("You may now inspect it with the 'inspect' command.\n")
		// Add it to the Pokedex and tag it for removal
		pokedex[pokemonName] = pokemon
		remove = true
	} else {
		fmt.Printf("%v %v!\n", pokemonName, color.YellowString("RESISTED"))

		if ranAway {
			// Tag it for removal
			fmt.Printf("Oh no, %v %v!\n", pokemonName, color.RedString("ran away"))
			remove = true
		} else {
			fmt.Println()
		}
	}
	if remove {
		// Remove it from the explored area
		if removedPokemons[currentLocationArea] == nil {
			removedPokemons[currentLocationArea] = map[string]bool{}
		}
		removedPokemons[currentLocationArea][pokemonName] = true
		delete(currentPokemons, pokemonName)
		time.Sleep(time.Duration(1) * time.Second)
		commandExplore(currentLocationArea)
	}

	return nil
}

/*
commandInspect
Inspects the given pokemon name if it exists in the pokedex.
*/
func commandInspect(name string) error {
	if name == "" {
		fmt.Printf("\nPlease provide the name of the Pokemon to inspect.\n")
		return nil
	}
	pokemon, ok := pokedex[name]
	if !ok {
		fmt.Printf("\n%v was not caught yet!\n", name)
		return nil
	}

	// Print basic info
	fmt.Printf("\nName: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)

	// Print Stats
	fmt.Println("Stats:")
	for i := range pokemon.Stats {
		baseStat := pokemon.Stats[i].BaseStat
		statName := pokemon.Stats[i].Stat.Name
		fmt.Printf("  - %v: %v\n", baseStat, statName)
	}

	// Print Types
	fmt.Println("Types:")
	for i := range pokemon.Types {
		fmt.Printf("  - %v\n", pokemon.Types[i].Type.Name)
	}

	return nil
}

/*
commandPokedex
Inspects the given pokemon name if it exists in the pokedex.
*/
func commandPokedex(_ string) error {
	if len(pokedex) < 1 {
		fmt.Printf("\nYou haven't caught any Pokemon yet!\n")
		return nil
	}

	fmt.Printf("\nPokemons you caught:\n")
	for i := range pokedex {
		fmt.Printf("  - %v\n", pokedex[i].Name)
	}

	return nil
}
