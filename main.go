package main

import (
	"bufio"
	"fmt"
	"gokedex/gokeapi"
	"math/rand"
	"os"
	"strings"
)

var (
	gokedex = map[string]gokeapi.PokemonInfo{}
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

func allCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Show next 20 areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show last 20 areas",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Show info about a specific area",
			callback:    exploreCommand,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback:    catchCommand,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon",
			callback:    inspectCommand,
		},
		"gokedex": {
			name:        "gokedex",
			description: "Show gokedex entries",
			callback:    gokedexCommand,
		},
	}
}

func commandHelp(_ []string) error {
	fmt.Println(`

Welcome to the Gokedex!
Usage:`)
	fmt.Println()
	for _, cmd := range allCommands() {
		fmt.Println(fmt.Sprintf("%s: %s", cmd.name, cmd.description))
	}
	return nil
}

func commandExit(_ []string) error {
	os.Exit(0)
	return nil
}

func commandMap(_ []string) error {
	areas, err := gokeapi.LoadNextAreas()
	if err != nil {
		return err
	}
	for _, area := range areas {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapB(_ []string) error {
	areas, err := gokeapi.LoadPreviousAreas()
	if err != nil {
		return err
	}
	for _, area := range areas {
		fmt.Println(area.Name)
	}
	return nil
}

func exploreCommand(args []string) error {
	if len(args) != 1 {
		fmt.Println("Expected exactly one argument")
		return nil
	}
	info, err := gokeapi.LoadAreaInfo(args[0])
	if err != nil {
		return err
	}
	fmt.Println("Found PokemonSummary:")
	for _, p := range info.PokemonEncounters {
		fmt.Println(fmt.Sprintf("- %s", p.Pokemon.Name))
	}
	return nil
}

func catchCommand(args []string) error {
	if len(args) != 1 {
		fmt.Println("Expected exactly one argument")
		return nil
	}
	name := args[0]
	fmt.Println(fmt.Sprintf("Throwing a Pokeball at %s...", name))
	pokemon, err := gokeapi.LoadPokemonInfo(name)
	if err != nil {
		return err
	}
	chance := float64(pokemon.BaseExp) / 255
	luck := rand.Float64()
	catch := luck < chance
	if catch {
		gokedex[name] = pokemon
		fmt.Println(fmt.Sprintf("%s was caught!", name))
		return nil
	}
	fmt.Println(fmt.Sprintf("%s escaped!", name))
	return nil
}

func inspectCommand(args []string) error {
	if len(args) != 1 {
		fmt.Println("Expected exactly one argument")
		return nil
	}
	name := args[0]
	pokemon, ok := gokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Println(fmt.Sprintf("Name: %s", pokemon.Name))
	fmt.Println(fmt.Sprintf("Height: %d", pokemon.Height))
	fmt.Println(fmt.Sprintf("Weight: %d", pokemon.Weight))
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Println(fmt.Sprintf("- %s: %d", stat.Stat.Name, stat.BaseStat))
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Println(fmt.Sprintf("- %s", t.Type.Name))
	}
	return nil
}

func gokedexCommand(args []string) error {
	fmt.Println("Your Gokedex:")
	for name, _ := range gokedex {
		fmt.Println(fmt.Sprintf("- %s", name))
	}
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Gokedex > ")
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		parts := strings.Split(text, " ")
		cmd, ok := allCommands()[parts[0]]
		if !ok {
			fmt.Println("Unknown command: ", text)
			fmt.Print("Gokedex > ")
			continue
		}
		if err := cmd.callback(parts[1:]); err != nil {
			fmt.Println(fmt.Sprintf("Error: %s", err.Error()))
			break
		}
		fmt.Print("Gokedex > ")
	}
}
