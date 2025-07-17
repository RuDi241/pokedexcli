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
		if len(words) > 0 {
			fmt.Printf("Your command was: %v\n", words[0])
		}
	}
}
