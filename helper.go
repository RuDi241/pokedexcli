package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/RuDi241/pokedexcli/pokecache"
)

var cache *pokecache.Cache

// Helper functions
func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

func getData[T any](url string) (data T, err error) {
	var responseData T
	// Check cache
	body, ok := cache.Get(url)
	if ok {
		json.Unmarshal(body, &responseData)
		return responseData, nil

	}
	res, err := http.Get(url)
	if err != nil {
		return responseData, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return responseData, fmt.Errorf("status code %v", res.StatusCode)
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return responseData, err
	}

	cache.Add(url, body)

	json.Unmarshal(body, &responseData)
	return responseData, nil
}
