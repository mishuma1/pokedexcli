package main

type Poke_location struct {
	Count    int
	Next     string
	Previous string
	Explore  string
	Results  []struct {
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
