package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
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

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
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

func prettyPrintEvent(event Event) {
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

func main() {
	events := getEvents()

	for i := 0; i < len(events); i++ {
		prettyPrintEvent(events[i])
	}

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	gcalEvents, err := srv.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(gcalEvents.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range gcalEvents.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
}
