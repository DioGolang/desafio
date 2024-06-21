package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	vars := mux.Vars(r)
	idStr := vars["eventID"]
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
	vars := mux.Vars(r)
	eventIDStr := vars["eventID"]
	log.Printf("Received request for eventID: %s", eventIDStr)

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		log.Printf("Invalid event ID: %s", eventIDStr)
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var eventSpots []Spot
	for _, spot := range data.Spots {
		if spot.EventID == eventID {
			eventSpots = append(eventSpots, spot)
		}
	}

	if len(eventSpots) == 0 {
		log.Printf("No spots found for event ID: %d", eventID)
		http.Error(w, "No spots found for this event", http.StatusNotFound)
		return
	}

	log.Printf("Found %d spots for event ID: %d", len(eventSpots), eventID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eventSpots)
}

func reserveSpot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventIDStr := vars["eventID"]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var spot Spot
	if err := json.NewDecoder(r.Body).Decode(&spot); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if spot.Name == "" {
		http.Error(w, "Spot name is required", http.StatusBadRequest)
		return
	}

	spot.ID = len(data.Spots) + 1
	spot.Status = "reserved"
	spot.EventID = eventID
	data.Spots = append(data.Spots, spot)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(spot)
}

func main() {
	loadData()

	r := mux.NewRouter()
	r.HandleFunc("/events", getEvents).Methods("GET")
	r.HandleFunc("/events/{eventID}", getEventByID).Methods("GET")
	r.HandleFunc("/events/{eventID}/spots", getSpotsByEventID).Methods("GET")
	r.HandleFunc("/events/{eventID}/reserve", reserveSpot).Methods("POST")

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
