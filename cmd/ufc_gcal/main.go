package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Event struct {
	Title  string  `json:"title"`
	Date   string  `json:"date"`
	Link   string  `json:"link"`
	Fights []Fight `json:"fights"`
}

type Fight struct {
	Main     bool    `json:"main"`
	Weight   string  `json:"weight"`
	FighterA Fighter `json:"fighter_a"`
	FighterB Fighter `json:"fighter_b"`
}

type Fighter struct {
	Name   string `json:"name"`
	Record string `json:"record"`
}

type MMAFightCardAPIResponse struct {
	ID        string  `json:"id"`
	Data      []Event `json:"data"`
	UpdatedAt string  `json:"updated_at"`
}

func getEvents() []Event {
	res, err := http.Get("https://mmafightcardsapi.adaptable.app/")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	var apiRes MMAFightCardAPIResponse

	err = json.Unmarshal(body, &apiRes)
	if err != nil {
		log.Fatal(err)
	}

	return apiRes.Data
}

func main() {
	events := getEvents()

	fmt.Println(events[0])
}
