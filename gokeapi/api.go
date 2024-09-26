package gokeapi

import (
	"encoding/json"
	"gokedex/gokecache"
	"io"
	"net/http"
	"time"
)

const (
	areasBaseUrl   = "https://pokeapi.co/api/v2/location-area/"
	pokemonBaseUrl = "https://pokeapi.co/api/v2/pokemon/"
	cacheInterval  = time.Minute
)

var (
	nextUrl     = areasBaseUrl
	previousUrl = areasBaseUrl
	cache       = gokecache.NewCache(cacheInterval)
)

type areaListResponse struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []AreaSummary `json:"results"`
}

type AreaSummary struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type AreaInfo struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon PokemonSummary `json:"pokemon"`
}

type PokemonSummary struct {
	Name string `json:"name"`
}

type PokemonInfo struct {
	Name    string `json:"name"`
	BaseExp int    `json:"base_experience"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
	Stats   []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func LoadNextAreas() ([]AreaSummary, error) {
	return loadAreas(nextUrl)
}

func LoadPreviousAreas() ([]AreaSummary, error) {
	return loadAreas(previousUrl)
}

func LoadAreaInfo(name string) (AreaInfo, error) {
	url := areasBaseUrl + name
	data, err := requestCached(url)
	if err != nil {
		return AreaInfo{}, err
	}
	var res AreaInfo
	err = json.Unmarshal(data, &res)
	return res, nil
}

func LoadPokemonInfo(name string) (PokemonInfo, error) {
	url := pokemonBaseUrl + name
	data, err := requestCached(url)
	if err != nil {
		return PokemonInfo{}, err
	}
	var res PokemonInfo
	err = json.Unmarshal(data, &res)
	return res, nil
}

func loadAreas(url string) ([]AreaSummary, error) {
	data, err := requestCached(url)
	if err != nil {
		return []AreaSummary{}, err
	}
	var res areaListResponse
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []AreaSummary{}, err
	}
	nextUrl = res.Next
	previousUrl = res.Previous
	return res.Results, nil
}

func requestCached(url string) ([]byte, error) {
	if val, ok := cache.Get(url); ok {
		return val, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	data, err := io.ReadAll(res.Body)
	return data, err
}
