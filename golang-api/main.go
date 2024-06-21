package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Organization string  `json:"organization"`
	Date         string  `json:"date"`
	Price        float64 `json:"price"`
	Rating       string  `json:"rating"`
	ImageURL     string  `json:"image_url"`
	CreatedAt    string  `json:"created_at"`
	Location     string  `json:"location"`
}

type Spot struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	EventID int    `json:"event_id"`
}

type Data struct {
	Events []Event `json:"events"`
	Spots  []Spot  `json:"spots"`
}

var data Data

const timeLayout = "2006-01-02T15:04:05"

func parseCustomTime(value string) (time.Time, error) {
	return time.Parse(timeLayout, value)
}

func loadData() {
	file, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Failed to read data file: %v", err)
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatalf("Failed to unmarshal data: %v", err)
	}

	// Convert string dates to time.Time
	for i, event := range data.Events {
		parsedDate, err := parseCustomTime(event.Date)
		if err != nil {
			log.Fatalf("Failed to parse event date: %v", err)
		}
		data.Events[i].Date = parsedDate.Format(time.RFC3339)

		parsedCreatedAt, err := parseCustomTime(event.CreatedAt)
		if err != nil {
			log.Fatalf("Failed to parse event created_at: %v", err)
		}
		data.Events[i].CreatedAt = parsedCreatedAt.Format(time.RFC3339)
	}
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data.Events)
}

func getEventByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/events/")
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	for _, event := range data.Events {
		if event.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(event)
			return
		}
	}
	http.Error(w, "Event not found", http.StatusNotFound)
}

func getSpotsByEventID(w http.ResponseWriter, r *http.Request) {
	// Extract event ID from the URL path using regular expression
	re := regexp.MustCompile(`/events/(\d+)/spots`)
	match := re.FindStringSubmatch(r.URL.Path)
	if len(match) != 2 {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
		return
	}
	idStr := match[1]

	// Convert event ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	// Filter spots based on the extracted event ID
	var eventSpots []Spot
	for _, spot := range data.Spots {
		if spot.EventID == id {
			eventSpots = append(eventSpots, spot)
		}
	}

	// Handle success and error scenarios
	if len(eventSpots) > 0 {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(eventSpots)
		if err != nil {
			return // Handle potential encoding error
		}
	} else {
		http.Error(w, "No spots found for this event", http.StatusNotFound)
	}
}

func reserveSpot(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) != 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(urlParts[2])
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	spotName := urlParts[3]

	for i, spot := range data.Spots {
		if spot.EventID == eventID && spot.Name == spotName && spot.Status == "available" {
			data.Spots[i].Status = "reserved"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Spot reserved"})
			return
		}
	}

	http.Error(w, "Spot not available", http.StatusNotFound)
}

func main() {
	loadData()

	// Registre cada função handler com padrões de rota específicos
	http.HandleFunc("/events", getEvents)
	http.HandleFunc("/events/", getEventByID)
	http.HandleFunc("/events/:id/spots", getSpotsByEventID)
	http.HandleFunc("/events/reserve/", reserveSpot)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
