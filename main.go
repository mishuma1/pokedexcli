package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mishuma1/pokemon/cache"
)

var LocationData Poke_location
var newCache cache.Cache
var MyPokemons map[string]Pokemon

func cleanInput(text string) []string {
	word_split := strings.Split(text, " ")
	ret_slice := []string{}
	for _, item := range word_split {
		if item != "" {
			t_item := strings.ToLower(item)
			ret_slice = append(ret_slice, t_item)
		}
	}
	return ret_slice
}

func main() {
	LocationData.Next = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	MyPokemons = map[string]Pokemon{}
	Commands = map[string]CliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    CommandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays 20 locations forward",
			callback:    CommandMapNext,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays 20 locations previous",
			callback:    CommandMapPrevious,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location for pokemon",
			callback:    CommandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try and catch a pokemon",
			callback:    CommandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "View a captured pokemon",
			callback:    CommandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List your captured pokemon",
			callback:    CommandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    CommandExit,
		},
	}

	newCache = cache.Cache{}
	newCache.NewCache(30, 300)

	for {
		input := bufio.NewScanner(os.Stdin)
		fmt.Print("Pokedex > ")
		for input.Scan() {

			text := input.Text()
			cleanWords := cleanInput(text)
			if len(cleanWords) >= 1 {
				commandToExecute, ok := Commands[cleanWords[0]]
				LocationData.FullCommand = cleanWords

				if !ok {
					fmt.Printf("Unknown command\n")
				} else {
					errCom := commandToExecute.callback(LocationData)
					if errCom != nil {
						fmt.Printf("%v\n", errCom.Error())
					}
				}

			}

			//fmt.Printf("Your command was: %v\n", cleanWords[0])
			fmt.Print("Pokedex > ")
		}
		if err := input.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
	}
}
