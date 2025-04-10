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

type Poke_location struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		Url  string
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(Poke_location) error
}

var Commands map[string]cliCommand
var LocationData Poke_location

func cleanInput(text string) []string {
	word_split := strings.Split(text, " ")
	ret_slice := []string{}
	for _, item := range word_split {
		//		fmt.Printf("text: %v, item: %v\n", text, item)
		if item != "" {
			t_item := strings.ToLower(item)
			ret_slice = append(ret_slice, t_item)
		}
	}

	//	fmt.Printf("WORDS LEN: %d, %v\n", len(ret_slice), ret_slice)
	return ret_slice
}

func commandExit(config Poke_location) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(config Poke_location) error {
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage:\n\n")

	for _, key := range Commands {
		fmt.Printf("%v: %v\n", key.name, key.description)
	}

	return nil
}

func commandMap(url string) error {
	fmt.Printf("Do map stuff\n")
	resp, errResp := http.Get(url)
	if errResp != nil {
		return fmt.Errorf("error getting pokemon location: %v", errResp)
	}

	jsonBytes, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return fmt.Errorf("error parsing poke data: %v", errRead)
	}

	errUn := json.Unmarshal(jsonBytes, &LocationData)
	if errUn != nil {
		return fmt.Errorf("error unmarshalling poke data: %v", errUn)
	}

	for _, val := range LocationData.Results {
		fmt.Printf("%v\n", val.Name)
	}

	return nil
}

func commandMapNext(config Poke_location) error {
	if config.Next == "" {
		config.Next = config.Previous
	}
	err := commandMap(config.Next)
	return err
}

func commandMapPrevious(config Poke_location) error {
	if config.Previous == "" {
		config.Previous = config.Next
	}
	err := commandMap(config.Previous)
	return err
}

func main() {
	LocationData.Next = "https://pokeapi.co/api/v2/location-area"
	Commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays 20 locations forward",
			callback:    commandMapNext,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays 20 locations previous",
			callback:    commandMapPrevious,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	for {
		//map config : https://pokeapi.co/api/v2/location-area
		//?offset=20&limit=20

		input := bufio.NewScanner(os.Stdin)
		fmt.Print("Pokedex > ")
		for input.Scan() {

			text := input.Text()
			cleanWords := cleanInput(text)
			commandToExecute, ok := Commands[cleanWords[0]]
			if !ok {
				fmt.Printf("Unknown command\n")
			} else {
				commandToExecute.callback(LocationData)
			}

			//fmt.Printf("Your command was: %v\n", cleanWords[0])
			fmt.Print("Pokedex > ")
		}
		if err := input.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
	}
}

// {
// 	"count": 1089,
// 	"next": "https://pokeapi.co/api/v2/location-area?offset=60&limit=20",
// 	"previous": "https://pokeapi.co/api/v2/location-area?offset=20&limit=20",
// 	"results": [
// 	  {
// 		"name": "solaceon-ruins-b3f-d",
// 		"url": "https://pokeapi.co/api/v2/location-area/41/"
// 	  },
// 	  {
// 		"name": "solaceon-ruins-b3f-e",
// 		"url": "https://pokeapi.co/api/v2/location-area/42/"
// 	  }
// 	]
//   }
