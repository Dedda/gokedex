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
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
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
		fmt.Printf("- %s\n", p.Pokemon.Name)
	}
	return nil
}

func catchCommand(args []string) error {
	if len(args) != 1 {
		fmt.Println("Expected exactly one argument")
		return nil
	}
	name := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	pokemon, err := gokeapi.LoadPokemonInfo(name)
	if err != nil {
		return err
	}
	chance := float64(pokemon.BaseExp) / 255
	luck := rand.Float64()
	catch := luck < chance
	if catch {
		gokedex[name] = pokemon
		fmt.Printf("%s was caught!\n", name)
		return nil
	}
	fmt.Printf("%s escaped!\n", name)
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
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("- %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("- %s\n", t.Type.Name)
	}
	return nil
}

func gokedexCommand(args []string) error {
	fmt.Println("Your Gokedex:")
	for name, _ := range gokedex {
		fmt.Printf("- %s", name)
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
			fmt.Printf("Error: %s\n", err.Error())
			break
		}
		fmt.Print("Gokedex > ")
	}
}
