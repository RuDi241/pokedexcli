package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

// Supported CLI commands
var commands map[string]cliCommand

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
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	message := "Welcome to the Pokedex!\nUsage:\n\n"

	for _, c := range commands {
		message += fmt.Sprintf("%v: %v\n", c.name, c.description)
	}

	fmt.Print(message)
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

			err := val.callback()
			if err != nil {
				fmt.Println("Error encountered")
			}
		} else {
			fmt.Println("Command not found")
		}

	}
}
