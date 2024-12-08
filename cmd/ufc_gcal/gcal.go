package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendar struct {
	service *calendar.Service
	ctx     context.Context
}

func NewGoogleCalendar() GoogleCalendar {
	ctx := context.Background()
	srv := getCalendarService(&ctx)

	return GoogleCalendar{
		service: srv,
		ctx:     ctx,
	}
}

func (gCal GoogleCalendar) GetCalendarEvents() *calendar.Events {
	t := time.Now().Format(time.RFC3339)
	gcalEvents, err := gCal.service.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	return gcalEvents
}

func (gCal GoogleCalendar) ShowCalendarEvents() {
	gcalEvents := gCal.GetCalendarEvents()
	prettyPrintCalendarEvents(gcalEvents)
}

func (gCal GoogleCalendar) AddCalendarEvent(title, description string, start, end time.Time) {
	ufcCalId := gCal.getCalendarIdByTitle("UFC")
	fmt.Println(ufcCalId)

	spew.Dump(start.String(), end)

	event := &calendar.Event{
		Summary:     title,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: "Asia/Kolkata",
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: "Asia/Kolkata",
		},
	}

	publishedEvent, err := gCal.service.Events.Insert(ufcCalId, event).Do()
	if err != nil {
		spew.Dump(err)
		log.Fatalf("Unable to create event. %v\n", err)
	}

	fmt.Printf("Event created: %s\n", publishedEvent.HtmlLink)
}

func (gCal GoogleCalendar) getCalendarIdByTitle(title string) string {
	listRes, err := gCal.service.CalendarList.List().Fields("items/id", "items/summary").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of calendars: %v", err)
	}

	var ufcCalId string
	for _, v := range listRes.Items {
		if v.Summary == title {
			ufcCalId = v.Id
		}
	}

	return ufcCalId
}

func prettyPrintCalendarEvents(gcalEvents *calendar.Events) {
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

func getCalendarService(ctx *context.Context) *calendar.Service {
	b, err := os.ReadFile("config/credentials/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(*ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	return srv
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "config/credentials/token.json"
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
