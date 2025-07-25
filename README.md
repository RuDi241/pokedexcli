# 🧭 PokedexCLI

A terminal-based REPL (Read-Eval-Print Loop) application that allows you to explore, catch, and inspect Pokémon from the PokéAPI.  
Navigate different location areas, view wild Pokémon, and build your personal Pokédex — all from the command line.

---

## 🚀 Features

- Browse location areas using pagination
- Explore wild Pokémon in specific areas
- Attempt to catch Pokémon with a probability-based system
- Inspect caught Pokémon details (stats, types, size)
- Maintain a persistent session Pokédex
- Built-in help and graceful exit

---

## 🛠 Installation

```bash
git clone https://github.com/RuDi241/pokedexcli.git
cd pokedexcli
go build -o pokedex
./pokedex
```

> ⚠ Requires Go 1.18+ installed on your system

---

## 💡 Usage

When launched, PokedexCLI runs in a REPL loop where you can enter the following commands:

### 🧾 Available Commands

| Command          | Description                                                      |
| ---------------- | ---------------------------------------------------------------- |
| `help`           | Displays this help message                                       |
| `exit`           | Exits the application                                            |
| `map`            | Displays the next 20 Pokémon location areas                      |
| `mapb`           | Displays the previous 20 Pokémon location areas                  |
| `explore <area>` | Lists Pokémon that can be encountered in the given location area |
| `catch <name>`   | Attempts to catch the specified Pokémon                          |
| `inspect <name>` | Displays detailed info about a caught Pokémon                    |
| `pokedex`        | Lists all Pokémon currently caught in your Pokédex               |

---

## 🧮 Catch Probability

The chance to catch a Pokémon is based on the following formula:

```
catch probability = (pokemonCatchChanceMultiplier) / log(baseExperience + 2)
```

Default value:

```
pokemonCatchChanceMultiplier = 2.5
```

If the probability exceeds a random threshold, the Pokémon is caught and added to your Pokédex.

---

## 🧪 Example Session

```
go run .
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
sunyshore-city-area
sinnoh-pokemon-league-area
oreburgh-mine-1f
oreburgh-mine-b1f
valley-windworks-area
eterna-forest-area
fuego-ironworks-area
mt-coronet-1f-route-207
mt-coronet-2f
mt-coronet-3f
mt-coronet-exterior-snowfall
mt-coronet-exterior-blizzard
mt-coronet-4f
mt-coronet-4f-small-room
mt-coronet-5f
mt-coronet-6f
mt-coronet-1f-from-exterior
Pokedex > explore canalave-city-area
Exploring canalave-city-area
Found Pokemon:
 - tentacool
 - tentacruel
 - staryu
 - magikarp
 - gyarados
 - wingull
 - pelipper
 - shellos
 - gastrodon
 - finneon
 - lumineon
Pokedex > catch tentacool
Throwing a Pokeball at tentacool...
Catch probability: 0.5904433431907248
tentacool was caught!
Pokedex > inspect tentacool
Name: tentacool
Height: 9
Weight: 455
Stats:
 -hp: 40
 -attack: 40
 -defense: 35
 -special-attack: 50
 -special-defense: 100
 -speed: 70
Types:
 - water
 - poison
Pokedex > pokedex
Your Pokedex:
 - tentacool
Pokedex > exit
Closing the Pokedex... Goodbye!
```
