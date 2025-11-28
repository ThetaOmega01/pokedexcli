# pokedexcli

A tiny interactive Pokédex you run in the terminal. It talks to the public PokeAPI to explore location areas, find Pokémon, attempt catches, and inspect your collection.

## Install

Requirements: Latest version of Go

```bash
# Clone and build
git clone https://github.com/ThetaOmega01/pokedexcli.git
cd pokedexcli
go build

# Or run without building
go run .
```

## Usage

Start the app, then type commands at the prompt:

```text
Pokedex > help
```

Commands:

- help — Show available commands
- exit — Quit the app
- map — Show the next 20 location areas (pagination forward)
- mapb — Show the previous 20 location areas (pagination back)
- explore <area_name> — List Pokémon that can appear in a location area
- catch <pokemon_name> — Throw a Poké Ball and try to catch a Pokémon
- inspect <pokemon_name> — Show details (stats, types, height/weight) of a caught Pokémon
- pokedex — List all Pokémon you have caught

### Examples

```text
Pokedex > map
Pokedex > explore canalave-city-area
Pokedex > catch pikachu
Pokedex > pokedex
Pokedex > inspect pikachu
```

## Notes

- Data source: https://pokeapi.co/
- Simple in-memory cache (TTL ~1 minute) reduces repeat API calls.
- Catch chance scales down as a Pokémon’s base experience increases; some will escape!

## Development

Run tests:

```bash
go test ./...
```

Format and tidy (optional):

```bash
go fmt ./...
go mod tidy
```
