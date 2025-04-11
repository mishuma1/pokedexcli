package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	fmt.Printf("MAP Trying: %s\n", config.Next)
	err := CommandMap(config.Next)
	return err
}

func CommandMapPrevious(config Poke_location) error {
	if config.Previous == "" {
		config.Previous = config.Next
	}
	fmt.Printf("MAPB Trying: %s\n", config.Previous)
	err := CommandMap(config.Previous)
	return err
}

func CommandExplore(config Poke_location) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", config.Explore)
	//jsonCache := newCache.GetCache(url)
	//if len(jsonCache) == 0 {
	resp, errResp := http.Get(url)
	if errResp != nil {
		return fmt.Errorf("error exploring pokemon location: %v", errResp)
	}

	jsonBytes, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return fmt.Errorf("error parsing poke explore data: %v", errRead)
	}
	//jsonCache = jsonBytes
	//newCache.AddCache(url, jsonCache)
	//} else {
	//	fmt.Println("Using Cache Data")
	//}
	//fmt.Printf("Explore JSON: %v", string(jsonBytes))

	ExploreData := List_Pokemon{}
	errUn := json.Unmarshal(jsonBytes, &ExploreData)
	if errUn != nil {
		return fmt.Errorf("error unmarshalling poke data: %v", errUn)
	}

	for _, val := range ExploreData.Pokemon_encounters {
		fmt.Printf("%v\n", val.Pokemon.Name)
	}

	return nil
}

//https://pokeapi.co/api/v2/location-area/
