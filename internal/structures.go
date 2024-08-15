package internal

type LocationAreas struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []LocationAreaSimple
}

type LocationAreaSimple struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationArea struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
	Name              string              `json:"name"`
}

type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name           string  `json:"name"`
	BaseExperience int     `json:"base_experience"`
	Height         int     `json:"height"`
	Weight         int     `json:"weight"`
	Stats          []Stats `json:"stats"`
	Types          []Types `json:"types"`
}

type Stats struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	}
}

type Types struct {
	Type struct {
		Name string `json:"name"`
	}
}
