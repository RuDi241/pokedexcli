package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/RuDi241/pokedexcli/pokecache"
)

// CONST
const (
	cacheLifetime = 10 * time.Second
	baseURL       = "https://pokeapi.co/api/v2/"
)

// Structs
type cliCommand struct {
	name        string
	description string
	callback    func(conf *commandData) error
}

type commandData struct {
	prevArea                     string
	nextArea                     string
	args                         []string
	pokemonCatchChanceMultiplier float64
}

// Supported CLI commands
var commands map[string]cliCommand

var cmdData commandData

var caughtPokemons map[string]pokemonData

func init() {
	commands = getCommands()
	cache = pokecache.NewCache(cacheLifetime)

	cmdData.prevArea = ""
	cmdData.pokemonCatchChanceMultiplier = 2.5
	cmdData.nextArea = baseURL + "location-area"
	caughtPokemons = make(map[string]pokemonData)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		ok := scanner.Scan()
		if !ok {
			break
		}
		text := scanner.Text()
		words := cleanInput(text)
		if len(words) <= 0 {
			fmt.Println("Enter a command")
		}
		if val, ok := commands[words[0]]; ok {

			cmdData.args = words[1:]
			err := val.callback(&cmdData)
			if err != nil {
				fmt.Printf("Error encountered: %v\n", err)
			}
		} else {
			fmt.Println("Command not found")
		}

	}
}
