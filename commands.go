package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

func getCommands() map[string]cliCommand {
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
			description: "Displays the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Takes <location area> as argument. Displays information about location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Takes <pokemon> as argument. Attempts to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Takes <pokemon> as argument. Inspect a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all caught pokemons",
			callback:    commandPokedex,
		},
	}
	return commands
}

func commandExit(conf *commandData) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *commandData) error {
	message := "Welcome to the Pokedex!\nUsage:\n\n"

	for _, c := range commands {
		message += fmt.Sprintf("%v: %v\n", c.name, c.description)
	}

	fmt.Print(message)
	return nil
}

func commandMap(conf *commandData) error {
	responseData, err := getData[locationAreaList](conf.nextArea)
	if err != nil {
		return err
	}

	conf.prevArea = conf.nextArea
	conf.nextArea = responseData.Next

	for _, la := range responseData.Results {
		println(la.Name)
	}

	return nil
}

func commandMapb(conf *commandData) error {
	if conf.prevArea == "" {
		println("you're on the first page")
		return nil
	}
	responseData, err := getData[locationAreaList](conf.nextArea)
	if err != nil {
		return err
	}
	conf.nextArea = conf.prevArea
	conf.prevArea = responseData.Prev

	for _, la := range responseData.Results {
		println(la.Name)
	}

	return nil
}

func commandExplore(conf *commandData) error {
	if len(conf.args) <= 0 {
		fmt.Println("Please provide location area as first argument")
		return nil
	}
	fullURL := baseURL + "location-area/" + conf.args[0]
	responseData, err := getData[locationAreaData](fullURL)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %v\nFound Pokemon:\n", conf.args[0])
	for _, r := range responseData.PokemonEncounters {
		fmt.Printf(" - %v\n", r.Pokemon.Name)
	}
	return nil
}

func commandCatch(conf *commandData) error {
	if len(conf.args) <= 0 {
		fmt.Println("Please provide pokemon name as first argument")
		return nil
	}
	fullURL := baseURL + "pokemon/" + conf.args[0]
	responseData, err := getData[pokemonData](fullURL)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", conf.args[0])
	catchProb := conf.pokemonCatchChanceMultiplier / math.Log(float64(responseData.BaseExperience)+2)
	fmt.Printf("Catch probability: %v \n", catchProb)

	if catchProb >= rand.Float64() {
		caughtPokemons[conf.args[0]] = responseData
		fmt.Printf("%v was caught!\n", conf.args[0])
	}

	return nil
}

func commandInspect(conf *commandData) error {
	if len(conf.args) <= 0 {
		fmt.Println("Please provide pokemon name as first argument")
		return nil
	}

	data, ok := caughtPokemons[conf.args[0]]
	if !ok {
		fmt.Printf("You haven't caught %v\n", conf.args[0])
		return nil
	}

	fmt.Printf("Name: %v\n", data.Name)
	fmt.Printf("Height: %v\n", data.Height)
	fmt.Printf("Weight: %v\n", data.Weight)
	fmt.Printf("Stats:\n")
	for _, s := range data.Stats {
		fmt.Printf("\t-%v: %v\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Printf("Types:\n")
	for _, t := range data.Types {
		fmt.Printf(" - %v\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(conf *commandData) error {
	fmt.Printf("Your Pokedex:\n")
	for k := range caughtPokemons {
		fmt.Printf(" - %v\n", k)
	}
	return nil
}
