package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func getLocationAreas(url string) (LocationAreas, error) {
	data, err := FetchData(url)
	if err != nil {
		return LocationAreas{}, err
	}

	var locationareas LocationAreas
	err = json.Unmarshal(*data, &locationareas)
	if err != nil {
		log.Fatal(err)
	}

	// Upate global url variables
	nextUrl = locationareas.Next
	previousUrl = locationareas.Previous

	return locationareas, nil
}

/*
printLocationAreas
Lists to the user the 20 next location areas to explore.
*/
func printLocationAreas(url string) error {
	locationareas, err := getLocationAreas(url)
	if err != nil {
		log.Fatal(err)
	}

	// Update current information
	currentLocationAreas = make(map[string]bool)
	currentLocationArea = ""
	currentPokemons = make(map[string]bool)

	fmt.Printf("\nCurrent availeble areas to explore:\n")
	for i := range locationareas.Results {
		currentLocationAreas[locationareas.Results[i].Name] = true
		fmt.Println(locationareas.Results[i].Name)
	}

	return nil
}

func FetchData(url string) (*[]byte, error) {
	// Check if  request is present in cache
	data, found := cache.Get(url)

	if !found {
		// Fetch data from url
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		*data, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Check if return is valid
		if string((*data)[:]) == "Not Found" {
			fmt.Printf("Incorrect name provided.\n")
			return &[]byte{}, errors.New("not found")
		}

		// Add entry to cache
		NewCache(&cache, url, data, time.Duration(cacheInterval)*time.Second)
	}

	return data, nil
}

func catchAttempt(pokemon Pokemon) (bool, bool) {
	catchChance := catchRate - pokemon.BaseExperience
	catchSuccessArray := make([]bool, catchRate)
	for i := 0; i < catchRate; i++ {
		if i <= catchChance {
			catchSuccessArray[i] = true
		} else {
			catchSuccessArray[i] = false
		}
	}
	success := catchSuccessArray[rand.Intn(len(catchSuccessArray))]
	ranAway := rand.Int()%3 == 0
	return success, ranAway
}
