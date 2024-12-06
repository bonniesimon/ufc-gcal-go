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
	FighterA Fighter `json:"fighterA"`
	FighterB Fighter `json:"fighterB"`
}

type Fighter struct {
	Name   string `json:"name"`
	Record string `json:"record"`
}

type MMAFightCardAPIResponse struct {
	ID        string  `json:"id"`
	Data      []Event `json:"data"`
	UpdatedAt string  `json:"updatedAt"`
}

func GetMMAEvents() []Event {
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

func PrettyPrintEvent(event Event) {
	istTime, err := convertETtoIST(event.Date)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n\n\n\n")
	fmt.Printf(`
		------------------------------- %s -------------------------------
		%v (Early prelim time)
		%s
	`, event.Title, istTime.Format("Monday, January 2, 3:04 PM MST"), event.Link)

	for i := 0; i < len(event.Fights); i++ {
		fmt.Printf(`
			%s vs %s
		`, event.Fights[i].FighterA.Name, event.Fights[i].FighterB.Name)
	}

	fmt.Println("------------------------------------------------------------------")
}
