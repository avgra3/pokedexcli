package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const pokiAPI = "https://pokeapi.co/api/v2/location-area"

func GetLocations(url string) LocationResult {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatal("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	locationResults := LocationResult{}
	err = json.Unmarshal(body, &locationResults)
	if err != nil {
		log.Fatal(err)
	}
	return locationResults
}

type Locations struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationResult struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []Locations `json:"results"`
}
