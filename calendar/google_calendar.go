package calendar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"

	cal "google.golang.org/api/calendar/v3"
)

var calendarClient *calendar.Service

// InitCalendarAPI ...
func InitCalendarAPI(credentialsLocation string, tokenFilename string) error {
	b, err := ioutil.ReadFile(credentialsLocation)
	if err != nil {
		log.Println(fmt.Sprintf("Error while reading file, filename: %s",credentialsLocation), err)

		return err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope)
	if err != nil {
		log.Println(fmt.Sprintf("Unable to parse client secret file to config"), err)

		return err
	}
	client := getClient(config, tokenFilename)

	calendarClient, err = calendar.New(client)
	if err != nil {
		log.Println(fmt.Sprintf("Unable to retrieve Calendar client"), err)

		return err
	}

	return nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokenFilename string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokenFilename)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFilename, tok)
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

// QueryReservationsBetweenDate ...
func QueryReservationsBetweenDate(fromDate string, toDate string, calendarID string) (*cal.Events, error) {
	log.Println("Query events from google calendar, calendarId: " + calendarID)
	events, err := calendarClient.Events.List(calendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(fromDate).TimeMax(toDate).OrderBy("startTime").Do()
	if err != nil {
		log.Println(fmt.Sprintf("Unable to retrieve events"), err)
		return nil, err
	}
	if len(events.Items) == 0 {
		log.Println("No upcoming events found.")
		return nil, nil
	}
	return events, nil
}

// DeleteEventByID ...
func DeleteEventByID(calendarID string, eventID string) error{
	err := calendarClient.Events.Delete(calendarID, eventID).Do()
	if err != nil {
		log.Println("Unable to delete event wit calendarId: " + calendarID + " eventID: " + eventID)
		return err
	}

	log.Println("Event deleted with event id: " + eventID)
	return nil
}

// GetEventDate ...
func GetEventDate(event *cal.Event) (from string, to string) {
	startDate := event.Start.DateTime
	if startDate == "" {
		startDate = event.Start.Date
	}
	endDate := event.End.DateTime
	if endDate == "" {
		endDate = event.End.Date
	}

	return startDate, endDate
}
