package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Helper functions
func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

// Structs
type cliCommand struct {
	name        string
	description string
	callback    func(conf *config) error
}

type config struct {
	prev string
	next string
}

type locationAreaResponse struct {
	Count   int            `json:"count"`
	Next    string         `json:"next"`
	Prev    string         `json:"previous"`
	Results []locationArea `json:"results"`
}

type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Supported CLI commands
var commands map[string]cliCommand

var conf config

func init() {
	commands = map[string]cliCommand{
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
	}

	conf.prev = ""
	conf.next = "https://pokeapi.co/api/v2/location-area"
}

func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config) error {
	message := "Welcome to the Pokedex!\nUsage:\n\n"

	for _, c := range commands {
		message += fmt.Sprintf("%v: %v\n", c.name, c.description)
	}

	fmt.Print(message)
	return nil
}

func commandMap(conf *config) error {
	res, err := http.Get(conf.next)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var responseData locationAreaResponse
	json.Unmarshal(body, &responseData)

	conf.prev = conf.next
	conf.next = responseData.Next

	for _, la := range responseData.Results {
		println(la.Name)
	}

	return nil
}

func commandMapb(conf *config) error {
	if conf.prev == "" {
		println("you're on the first page")
		return nil
	}
	res, err := http.Get(conf.prev)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var responseData locationAreaResponse
	json.Unmarshal(body, &responseData)

	conf.next = conf.prev
	conf.prev = responseData.Prev

	for _, la := range responseData.Results {
		println(la.Name)
	}

	return nil
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

			err := val.callback(&conf)
			if err != nil {
				fmt.Printf("Error encountered: %v\n", err)
			}
		} else {
			fmt.Println("Command not found")
		}

	}
}
