package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/RuDi241/pokedexcli/pokecache"
)

// CONST
const (
	cacheLifetime = 10 * time.Second
	baseURL       = "https://pokeapi.co/api/v2/"
)

var cache *pokecache.Cache

// Helper functions
func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

func getData[T any](url string) (data T, err error) {
	var responseData T
	// Check cache
	body, ok := cache.Get(url)
	if ok {
		json.Unmarshal(body, &responseData)
		return responseData, nil

	}
	res, err := http.Get(url)
	if err != nil {
		return responseData, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return responseData, fmt.Errorf("status code %v", res.StatusCode)
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return responseData, err
	}

	cache.Add(url, body)

	json.Unmarshal(body, &responseData)
	return responseData, nil
}

// Structs
type cliCommand struct {
	name        string
	description string
	callback    func(conf *commandData) error
}

type commandData struct {
	prevArea               string
	nextArea               string
	args                   []string
	basePokemonCatchChance float64
}

type locationAreaList struct {
	Count   int    `json:"count"`
	Next    string `json:"next"`
	Prev    string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type locationAreaData struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type pokemonData struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        interface{} `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  interface{} `json:"ability"`
			IsHidden bool        `json:"is_hidden"`
			Slot     int         `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastTypes []interface{} `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string      `json:"front_default"`
				FrontFemale  interface{} `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string      `json:"back_default"`
				BackFemale       string      `json:"back_female"`
				BackShiny        string      `json:"back_shiny"`
				BackShinyFemale  interface{} `json:"back_shiny_female"`
				FrontDefault     string      `json:"front_default"`
				FrontFemale      string      `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale string      `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       string `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      string `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  string `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

// Supported CLI commands
var commands map[string]cliCommand

var cmdData commandData

var caughtPokemons map[string]pokemonData

func init() {
	cache = pokecache.NewCache(cacheLifetime)
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

	cmdData.prevArea = ""
	cmdData.basePokemonCatchChance = 0.5
	cmdData.nextArea = baseURL + "location-area"
	caughtPokemons = make(map[string]pokemonData)
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
	catchChance := conf.basePokemonCatchChance / math.Log(float64(responseData.BaseExperience)+2)

	if catchChance >= rand.Float64() {
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
