package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
)

type CliCommand struct {
	name        string
	description string
	callback    func(Poke_location) error
}

var Commands map[string]CliCommand

func CommandExit(config Poke_location) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func CommandHelp(config Poke_location) error {
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage:\n\n")

	for _, key := range Commands {
		fmt.Printf("%v: %v\n", key.name, key.description)
	}

	return nil
}

func CommandMap(url string) error {
	jsonCache := newCache.GetCache(url)
	if len(jsonCache) == 0 {
		resp, errResp := http.Get(url)
		if errResp != nil {
			return fmt.Errorf("error getting pokemon location: %v", errResp)
		}

		jsonBytes, errRead := io.ReadAll(resp.Body)
		if errRead != nil {
			return fmt.Errorf("error parsing poke data: %v", errRead)
		}
		jsonCache = jsonBytes
		newCache.AddCache(url, jsonCache)
	} else {
		fmt.Println("Using Cache Data")
	}

	errUn := json.Unmarshal(jsonCache, &LocationData)
	if errUn != nil {
		return fmt.Errorf("error unmarshalling poke data: %v", errUn)
	}

	if LocationData.Previous == "" {
		LocationData.Previous = url
	}
	if LocationData.Next == "" {
		LocationData.Next = url
	}

	for _, val := range LocationData.Results {
		fmt.Printf("%v\n", val.Name)
	}

	return nil
}

func CommandMapNext(config Poke_location) error {
	err := CommandMap(config.Next)
	return err
}

func CommandMapPrevious(config Poke_location) error {
	if config.Previous == "" {
		config.Previous = config.Next
	}
	err := CommandMap(config.Previous)
	return err
}

func CommandExplore(config Poke_location) error {

	if len(config.FullCommand) < 2 {
		return fmt.Errorf("missing location")
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", config.FullCommand[1])

	jsonCache := newCache.GetCache(url)
	if len(jsonCache) == 0 {

		resp, errResp := http.Get(url)
		if errResp != nil {
			return fmt.Errorf("error exploring pokemon location: %v", errResp)
		}

		if resp.StatusCode == 404 {
			return fmt.Errorf("location data: %v not found", config.FullCommand[1])
		}

		jsonBytes, errRead := io.ReadAll(resp.Body)
		if errRead != nil {
			return fmt.Errorf("error parsing poke explore data: %v", errRead)
		}
		jsonCache = jsonBytes
		newCache.AddCache(url, jsonCache)
	} else {
		fmt.Println("Using Cache Data")
	}

	ExploreData := List_Pokemon{}
	errUn := json.Unmarshal(jsonCache, &ExploreData)
	if errUn != nil {
		return fmt.Errorf("explore location: not found")
	}

	for _, val := range ExploreData.Pokemon_encounters {
		fmt.Printf("%v\n", val.Pokemon.Name)
	}

	return nil
}

func CommandCatch(config Poke_location) error {

	if len(config.FullCommand) < 2 {
		return fmt.Errorf("missing pokemon")
	}

	_, ok := MyPokemons[config.FullCommand[1]]
	if ok {
		return fmt.Errorf("already caught a %v", config.FullCommand[1])
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", config.FullCommand[1])

	jsonCache := newCache.GetCache(url)
	if len(jsonCache) == 0 {

		resp, errResp := http.Get(url)
		if errResp != nil {
			return fmt.Errorf("error pulling pokemon species data: %v", errResp)
		}
		if resp.StatusCode == 404 {
			return fmt.Errorf("species data: %v not found", config.FullCommand[1])
		}

		jsonBytes, errRead := io.ReadAll(resp.Body)
		if errRead != nil {
			return fmt.Errorf("error parsing poke species data: %v", errRead)
		}
		jsonCache = jsonBytes
		newCache.AddCache(url, jsonCache)
	} else {
		fmt.Println("Using Cache Data")
	}

	SpeciesData := Pokemon{}
	errUn := json.Unmarshal(jsonCache, &SpeciesData)
	if errUn != nil {
		return fmt.Errorf("unmarshal error: %v", errUn.Error())
	}

	SpeciesData = AssignToFields(SpeciesData)
	fmt.Printf("Throwing a Pokeball at %v...\n", SpeciesData.Name)

	//Try and capture
	chance := rand.Intn(SpeciesData.Base_experience * 12)
	if chance/10 >= SpeciesData.Base_experience {
		fmt.Printf("%v was caught!\n", SpeciesData.Name)
		MyPokemons[SpeciesData.Name] = SpeciesData
		return nil
	}
	fmt.Printf("%v escaped!\n", SpeciesData.Name)

	return nil
}

func CommandInspect(config Poke_location) error {

	if len(config.FullCommand) < 2 {
		return fmt.Errorf("missing pokemon")
	}

	pokemon, ok := MyPokemons[config.FullCommand[1]]
	if !ok {
		return fmt.Errorf("you have not caught a %v", config.FullCommand[1])
	}
	PrintPokemonInfo(pokemon)
	return nil
}

func CommandPokedex(config Poke_location) error {
	fmt.Printf("Your Pokedex:\n")
	for idx := range MyPokemons {
		fmt.Printf("\t- %v\n", idx)
	}
	return nil
}
