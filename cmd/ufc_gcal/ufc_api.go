package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Event struct {
	Title      string    `json:"title"`
	DateString string    `json:"date"`
	Link       string    `json:"link"`
	Fights     []Fight   `json:"fights"`
	DateTime   time.Time // Saving as IST time for now
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

	// Process DateString into time.Time and update Event.DateTime
	for i, v := range apiRes.Data {
		istTime, _ := convertETtoIST(v.DateString)
		apiRes.Data[i].DateTime = istTime
	}

	return apiRes.Data
}

func PrettyPrintEvent(event Event) {
	fmt.Printf("\n\n\n\n")
	fmt.Printf(`
		------------------------------- %s -------------------------------
		%v (Early prelim time)
		%s
	`, event.Title, event.DateTime.Format("Monday, January 2, 3:04 PM MST"), event.Link)

	for i := 0; i < len(event.Fights); i++ {
		fmt.Printf(`
			%s vs %s
		`, event.Fights[i].FighterA.Name, event.Fights[i].FighterB.Name)
	}

	fmt.Println("------------------------------------------------------------------")
}

func GetEventAsString(event Event) string {
	result := ""
	header := fmt.Sprintf("%s\n%v (Early prelim time)\n%s\n", event.Title, event.DateTime.Format("Monday, January 2, 3:04 PM MST"), event.Link)
	result += header

	for i := 0; i < len(event.Fights); i++ {
		result += fmt.Sprintf("\n%s vs %s\n", event.Fights[i].FighterA.Name, event.Fights[i].FighterB.Name)
	}

	return result
}
