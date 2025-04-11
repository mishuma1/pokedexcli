package main

import "fmt"

type Poke_location struct {
	Count       int
	Next        string
	Previous    string
	FullCommand []string
	Results     []struct {
		Name string
		Url  string
	}
}

type List_Pokemon struct {
	Pokemon_encounters []struct {
		Pokemon struct {
			Name string
			Url  string
		}
	}
}

type Pokemon struct {
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Stats  []PokeStats
	Types  []PokeTypes
	//0-255 - higher the easiers
	Base_experience int
	Hp              int
	Attack          int
	Defense         int
	SpecialAttack   int
	SpecialDefense  int
	Speed           int
}

type PokeStats struct {
	Base_stat int
	Stat      struct {
		Name string
	}
}

type PokeTypes struct {
	Type struct {
		Name string
	}
}

// "types": [
//     {
//       "slot": 1,
//       "type": {
//         "name": "electric",
//         "url": "https://pokeapi.co/api/v2/type/13/"
//       }
//     }

func AssignToFields(poke Pokemon) Pokemon {
	for idx := range poke.Stats {
		// fmt.Printf("Stat: %v\n", poke.Stats[idx].Stat.Name)
		switch poke.Stats[idx].Stat.Name {
		case "hp":
			poke.Hp = poke.Stats[idx].Base_stat
		case "attack":
			poke.Attack = poke.Stats[idx].Base_stat
		case "defense":
			poke.Defense = poke.Stats[idx].Base_stat
		case "special-attack":
			poke.SpecialAttack = poke.Stats[idx].Base_stat
		case "special-defense":
			poke.SpecialDefense = poke.Stats[idx].Base_stat
		case "speed":
			poke.Speed = poke.Stats[idx].Base_stat
		}
	}

	return poke
}

func PrintPokemonInfo(poke Pokemon) Pokemon {
	fmt.Printf("Name: %v\n", poke.Name)
	fmt.Printf("Height: %v\n", poke.Height)
	fmt.Printf("Weight: %v\n", poke.Weight)
	fmt.Printf("Stats:\n")
	fmt.Printf("\t-hp: %v\n", poke.Hp)
	fmt.Printf("\t-attack: %v\n", poke.Attack)
	fmt.Printf("\t-defense: %v\n", poke.Defense)
	fmt.Printf("\t-special-attack: %v\n", poke.SpecialAttack)
	fmt.Printf("\t-special-defense: %v\n", poke.SpecialDefense)
	fmt.Printf("\t-speed: %v\n", poke.Speed)
	fmt.Printf("Types:\n")
	for idx2 := range poke.Types {
		fmt.Printf("\t-%v\n", poke.Types[idx2].Type.Name)
	}

	return poke
}
